import assert from 'node:assert/strict';
import { readFileSync } from 'node:fs';
import test from 'node:test';
import ts from 'typescript';

const source = readFileSync(new URL('../src/stores/queryTabsStore.ts', import.meta.url), 'utf8');
const ast = ts.createSourceFile('store.ts', source, ts.ScriptTarget.Latest, true);
let moveSource;
function visit(node) {
    if (ts.isFunctionDeclaration(node) && node.name?.text === 'moveTab') moveSource = node.getText(ast);
    ts.forEachChild(node, visit);
}
visit(ast);
assert.ok(moveSource, 'Store exposes a moveTab action');
const compiled = ts.transpileModule(moveSource, { compilerOptions: { target: ts.ScriptTarget.ES2022 } }).outputText;

test('tabs move in either direction while retaining their content and identity', () => {
    const originals = ['a', 'b', 'c', 'd'].map(id => ({ id, sql: id, loading: id === 'b' }));
    const tabs = { value: [...originals] };
    const move = new Function('tabs', `${compiled}; return moveTab;`)(tabs);
    move('a', 'c', 'after');
    assert.deepEqual(tabs.value.map(tab => tab.id), ['b', 'c', 'a', 'd']);
    move('d', 'b', 'before');
    assert.deepEqual(tabs.value.map(tab => tab.id), ['d', 'b', 'c', 'a']);
    move('a', 'd', 'before');
    move('b', 'c', 'before');
    assert.deepEqual(tabs.value.map(tab => tab.id), ['a', 'd', 'b', 'c']);
    for (const original of originals) assert.equal(tabs.value.find(tab => tab.id === original.id), original);
    const before = [...tabs.value];
    move('a', 'a', 'after');
    move('missing', 'b', 'before');
    move('a', 'missing', 'after');
    assert.deepEqual(tabs.value, before);
});
