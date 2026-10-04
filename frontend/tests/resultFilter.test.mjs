import assert from 'node:assert/strict';
import { createRequire } from 'node:module';
import { readFileSync } from 'node:fs';
import test from 'node:test';
import ts from 'typescript';

const source = readFileSync(new URL('../src/utils/resultFilter.ts', import.meta.url), 'utf8');
const compiled = ts.transpileModule(source, { compilerOptions: { target: ts.ScriptTarget.ES2022, module: ts.ModuleKind.CommonJS } }).outputText;
const module = { exports: {} };
new Function('module', 'exports', 'require', compiled)(module, module.exports, createRequire(import.meta.url));
const { compileResultFilter } = module.exports;
const columns = ['first_name', 'age', 'minimum_age'];
const rows = [['Penelope', 21, 18], ['Nick', 17, 18], ['Penelope', 16, 16]];

test('constant comparisons combine with column filters', () => {
    for (const [expression, expected] of [
        ["first_name like 'Penelope' and 1 = 0", [false, false, false]],
        ["first_name like 'Penelope' and 1 = 1", [true, false, true]],
        ["1 = 0 OR first_name = 'Nick'", [false, true, false]],
        ["(1 = 0 OR 2 > 1) AND age >= 18", [true, false, false]],
        ["'a' = 'a' AND TRUE = TRUE", [true, true, true]],
        ["18 <= age AND age >= minimum_age", [true, false, false]],
        ["NULL IS NULL", [true, true, true]],
        ["1 IN (1, 2) AND 'Penelope' LIKE 'P%'", [true, true, true]],
    ]) assert.deepEqual(rows.map(compileResultFilter(expression, columns)), expected, expression);
});

test('existing operators and validation remain supported', () => {
    assert.deepEqual(rows.map(compileResultFilter("first_name IN ('Penelope') AND age < 18", columns)), [false, false, true]);
    assert.deepEqual(rows.map(compileResultFilter("NOT (age < 18) OR first_name = 'Nick'", columns)), [true, true, false]);
    for (const expression of ['1 =', 'missing = 0', '1 = 0 garbage', "first_name = 'unterminated"]) {
        assert.throws(() => compileResultFilter(expression, columns), expression);
    }
});
