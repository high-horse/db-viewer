import assert from 'node:assert/strict';
import { readFileSync } from 'node:fs';
import { createRequire } from 'node:module';
import test from 'node:test';
import ts from 'typescript';
import { parse, compileScript } from '@vue/compiler-sfc';
import { createRenderer, nextTick } from 'vue';
import { createPinia, setActivePinia } from 'pinia';

const require = createRequire(import.meta.url);
let saveResponse = { applied: 0, error: '' };
const backend = { DescribeTableEdit: async () => null, SaveTableChanges: async () => saveResponse };
function evaluate(source, overrides = {}) {
    const output = ts.transpileModule(source, { compilerOptions: { module: ts.ModuleKind.CommonJS, target: ts.ScriptTarget.ES2022 } }).outputText;
    const module = { exports: {} };
    new Function('module', 'exports', 'require', output)(module, module.exports, name => overrides[name] ?? require(name));
    return module.exports;
}
const stores = evaluate(readFileSync(new URL('../src/stores/tableEditsStore.ts', import.meta.url), 'utf8'), { '@bindings/db-viewer/internal/app': { DbService: backend } });
const { descriptor } = parse(readFileSync(new URL('../src/components/workspace/resultGrid/TableDataEditor.vue', import.meta.url), 'utf8'));
const componentSource = compileScript(descriptor, { id: 'table-editor-test' }).content;
const component = evaluate(componentSource, {
    '@/stores/tableEditsStore': stores,
    '@/stores/queryTabsStore': { useQueryTabsStore: () => ({ tabs: [], execute: async () => {} }) },
    quasar: { Notify: { create() {} } },
}).default;
component.render = () => null;
const renderer = createRenderer({
    createComment: text => ({ text }), createText: text => ({ text }), createElement: tag => ({ tag }),
    insert() {}, remove() {}, setText() {}, setElementText() {}, patchProp() {}, parentNode() { return null; }, nextSibling() { return null; },
});
function setup(driver = 'sqlite') {
    setActivePinia(createPinia());
    const edits = stores.useTableEditsStore();
    edits.states.tab = { info: { driver, table: { connectionId: 'conn', schema: 'main', database: 'db', name: 'items' }, canInsert: true, canModify: true, keys: driver === 'mongodb' ? ['_id'] : ['id'], reason: '', columns: [
        { name: 'id', databaseType: 'INTEGER', nullable: false, autoIncrement: true, generated: false },
        { name: 'name', databaseType: 'TEXT', nullable: false, autoIncrement: false, generated: false },
        { name: 'note', databaseType: 'TEXT', nullable: true, defaultValue: 'default', autoIncrement: false, generated: false },
        { name: 'generated', databaseType: 'INTEGER', generated: true },
    ] }, drafts: [], saving: false, loading: false, error: '' };
    const result = { Columns: [{ Name: 'id' }, { Name: 'name' }, { Name: 'note' }, { Name: 'generated' }], Rows: [[1, 'before', null, 2]], StartRow: 1, PageSize: 100 };
    const app = renderer.createApp(component, { tabId: 'tab', result, selectedKey: '' });
    const vm = app.mount({});
    const setup = vm.$.setupState;
    return { edits, result, vm, setup, app };
}

test('SQL row editing stages changes, preserves null versus empty, and coalesces edits', async () => {
    const { edits, result, vm, setup, app } = setupEditor();
    const key = edits.rowKey('tab', result, 0);
    vm.openRow(key);
    await nextTick();
    setup.fields.find(field => field.name === 'name').text = 'edited';
    const note = setup.fields.find(field => field.name === 'note');
    note.mode = 'value'; note.text = '';
    setup.stageRow();
    assert.equal(edits.states.tab.drafts.length, 1);
    assert.deepEqual(edits.states.tab.drafts[0].values, { name: 'edited', note: '' });
    assert.deepEqual(edits.states.tab.drafts[0].keys, { id: 1 });
    assert.equal(result.Rows[0][1], 'before');
    vm.openRow(key);
    setup.fields.find(field => field.name === 'name').text = 'edited again';
    setup.stageRow();
    assert.equal(edits.states.tab.drafts.length, 1);
    assert.equal(edits.states.tab.drafts[0].values.name, 'edited again');
    edits.discard('tab');
    assert.equal(edits.hasPending('tab'), false);
    app.unmount();
});
const setupEditor = setup;

test('new SQL rows omit generated/default fields until saved', () => {
    const { edits, vm, setup, app } = setupEditor();
    vm.openRow();
    setup.fields.find(field => field.name === 'name').text = 'new';
    setup.stageRow();
    assert.equal(edits.states.tab.drafts[0].operation, 'insert');
    assert.deepEqual(edits.states.tab.drafts[0].values, { name: 'new' });
    assert.deepEqual(edits.states.tab.drafts[0].keys, {});
    app.unmount();
});

test('Mongo document staging preserves Extended JSON and absent keys', () => {
    const { edits, result, vm, setup, app } = setupEditor('mongodb');
    result.Documents = ['{"_id":{"$oid":"507f1f77bcf86cd799439011"},"large":{"$numberLong":"9007199254740993"}}'];
    const key = edits.rowKey('tab', result, 0);
    vm.openRow(key);
    const doc = JSON.parse(setup.documentText);
    assert.equal(Object.hasOwn(doc, '_id'), false);
    doc.name = 'edited';
    setup.documentText = JSON.stringify(doc);
    setup.stageRow();
    assert.deepEqual(edits.states.tab.drafts[0].values.large, { $numberLong: '9007199254740993' });
    assert.equal(Object.hasOwn(edits.states.tab.drafts[0].values, 'note'), false);
    assert.deepEqual(edits.states.tab.drafts[0].keys._id, JSON.parse(result.Documents[0])._id);
    vm.openRow(key);
    setup.documentText = '{"_id":"changed"}';
    setup.stageRow();
    assert.match(setup.dialogError, /_id is read-only/);
    assert.equal(edits.states.tab.drafts.length, 1);
    app.unmount();
});

test('a partially saved batch keeps only failed/unattempted drafts for retry', async () => {
    const { edits, app } = setupEditor('mongodb');
    edits.stage('tab', { key: 'one', operation: 'insert', keys: {}, values: { name: 'one' } });
    edits.stage('tab', { key: 'two', operation: 'insert', keys: {}, values: { name: 'two' } });
    saveResponse = { applied: 1, error: 'second insert failed' };
    assert.equal(await edits.save('tab'), true);
    assert.deepEqual(edits.states.tab.drafts.map(row => row.key), ['two']);
    assert.equal(edits.states.tab.error, 'second insert failed');
    saveResponse = { applied: 1, error: '' };
    assert.equal(await edits.save('tab'), true);
    assert.equal(edits.hasPending('tab'), false);
    app.unmount();
});


test('existing primary keys are locked and excluded from staged updates', () => {
    const { edits, result, vm, setup, app } = setupEditor();
    vm.openRow(edits.rowKey('tab', result, 0));
    const identity = setup.fields.find(field => field.name === 'id');
    assert.equal(identity.locked, true);
    identity.text = '999'; // Even a tampered UI field cannot stage an identity change.
    setup.fields.find(field => field.name === 'name').text = 'edited';
    setup.stageRow();
    assert.deepEqual(edits.states.tab.drafts[0].values, { name: 'edited' });
    assert.deepEqual(edits.states.tab.drafts[0].keys, { id: 1 });
    vm.openRow();
    assert.equal(setup.fields.find(field => field.name === 'id').locked, true);
    app.unmount();
});
