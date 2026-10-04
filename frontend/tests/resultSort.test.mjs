import assert from 'node:assert/strict';
import { readFileSync } from 'node:fs';
import test from 'node:test';
import ts from 'typescript';

function load(name) {
    const source = readFileSync(new URL(`../src/utils/${name}.ts`, import.meta.url), 'utf8');
    const compiled = ts.transpileModule(source, { compilerOptions: { target: ts.ScriptTarget.ES2022, module: ts.ModuleKind.CommonJS } }).outputText;
    const module = { exports: {} };
    new Function('module', 'exports', compiled)(module, module.exports);
    return module.exports;
}
const { sortResultRows } = load('resultSort');
const { compileResultFilter } = load('resultFilter');

test('descending sort retains the filtered records and their original row indexes', () => {
    const source = [['Penelope', 2], ['Nick', 99], ['Penelope', 10], ['Penelope', 1]];
    const matches = compileResultFilter("first_name LIKE 'Penelope' AND 1 = 1", ['first_name', 'id']);
    const filtered = source.map((row, index) => ({ row, index })).filter(({ row }) => matches(row));
    const descending = sortResultRows(filtered, 2, 'desc');
    assert.deepEqual(descending.map(({ row }) => row[1]), [10, 2, 1]);
    assert.deepEqual(descending.map(({ index }) => index), [2, 0, 3]);
    assert.deepEqual(sortResultRows(filtered, 2, 'asc').map(({ row }) => row[1]), [1, 2, 10]);
    assert.deepEqual(sortResultRows(filtered, 0), filtered);
    assert.deepEqual(filtered.map(({ row }) => row[1]), [2, 10, 1]);
    assert.deepEqual(source[1], ['Nick', 99]);
});

test('sorting handles nulls, stable ties, and BSON integers without precision loss', () => {
    const rows = [null, '{"$numberLong":"9007199254740993"}', '{"$numberLong":"9007199254740992"}'].map((value, index) => ({ row: [value], index }));
    assert.deepEqual(sortResultRows(rows, 1, 'desc').map(({ index }) => index), [1, 2, 0]);
    const ties = [{ row: ['same'], index: 1 }, { row: ['same'], index: 0 }];
    assert.deepEqual(sortResultRows(ties, 1, 'desc').map(({ index }) => index), [0, 1]);
});
