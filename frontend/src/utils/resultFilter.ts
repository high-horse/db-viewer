type Token = { text: string; kind: 'word' | 'string' | 'identifier' | 'number' | 'symbol' };
type Predicate = (row: unknown[]) => boolean;

// Parse a deliberately small expression language without evaluating user code.
export function compileResultFilter(expression: string, columns: string[]): Predicate {
    const tokens: Token[] = [];
    const pattern = /\s+|'(?:''|[^'])*'|"(?:""|[^"])*"|`(?:``|[^`])*`|(?:>=|<=|<>|!=|[=><(),])|-?\d+(?:\.\d+)?(?:[eE][+-]?\d+)?|[\p{L}_$][\p{L}\p{N}_$.]*/uy;
    let offset = 0;
    while (offset < expression.length) {
        pattern.lastIndex = offset;
        const match = pattern.exec(expression);
        if (!match) throw new Error(`Unexpected character at position ${offset + 1}`);
        offset = pattern.lastIndex;
        const text = match[0];
        if (/^\s/.test(text)) continue;
        tokens.push({ text, kind: text[0] === "'" ? 'string' : /^["`]/.test(text) ? 'identifier' : /^-?\d/.test(text) ? 'number' : /^[\p{L}_$]/u.test(text) ? 'word' : 'symbol' });
    }
    let position = 0;
    const peek = (word: string) => tokens[position]?.kind !== 'string' && tokens[position]?.kind !== 'identifier' && tokens[position]?.text.toUpperCase() === word;
    const take = (word: string) => { if (!peek(word)) return false; position++; return true; };
    const expect = (word: string) => { if (!take(word)) throw new Error(`Expected ${word}`); };
    const literal = (): unknown => {
        const token = tokens[position++];
        if (!token) throw new Error('Expected a value');
        if (token.kind === 'string') return token.text.slice(1, -1).replace(/''/g, "'");
        if (token.kind === 'number') return Number(token.text);
        if (token.kind === 'word') {
            if (token.text.toUpperCase() === 'TRUE') return true;
            if (token.text.toUpperCase() === 'FALSE') return false;
            if (token.text.toUpperCase() === 'NULL') return null;
        }
        throw new Error('Use single quotes for text values');
    };
    const normalize = (value: unknown): unknown => {
        if (typeof value !== 'string') return value;
        // MongoDB canonical JSON wraps BSON scalar values.
        try {
            const parsed = JSON.parse(value);
            if (parsed && typeof parsed === 'object') {
                if ('$oid' in parsed) return parsed.$oid;
                if ('$numberInt' in parsed) return Number(parsed.$numberInt);
                if ('$numberLong' in parsed) return BigInt(parsed.$numberLong);
                if ('$numberDouble' in parsed) return Number(parsed.$numberDouble);
            }
        } catch { /* Plain text remains plain text. */ }
        return value;
    };
    const compare = (raw: unknown, target: unknown, operator: string): boolean => {
        const value = normalize(raw);
        if (value == null || target == null) return false;
        let left: any = value, right: any = target;
        if (typeof left === 'bigint' && typeof right === 'number') {
            if (!Number.isSafeInteger(right)) throw new Error('Quote large integer values to preserve precision');
            right = BigInt(right);
        } else if (typeof left === 'bigint' && typeof right === 'string' && /^-?\d+$/.test(right)) right = BigInt(right);
        else if (typeof left !== typeof right) return false;
        switch (operator) {
            case '=': return left === right;
            case '!=': case '<>': return left !== right;
            case '>': return left > right;
            case '<': return left < right;
            case '>=': return left >= right;
            case '<=': return left <= right;
        }
        return false;
    };
    const operand = (): ((row: unknown[]) => unknown) => {
        const token = tokens[position];
        if (!token) throw new Error('Expected a column name or value');
        if (token.kind === 'number' || token.kind === 'string' ||
            (token.kind === 'word' && ['TRUE', 'FALSE', 'NULL'].includes(token.text.toUpperCase()))) {
            const value = literal();
            return () => value;
        }
        position++;
        if (!['word', 'identifier'].includes(token.kind)) throw new Error('Expected a column name or value');
        const name = token.kind === 'identifier' ? token.text.slice(1, -1).replaceAll(token.text[0]!.repeat(2), token.text[0]!) : token.text;
        let index = columns.indexOf(name);
        if (index < 0) index = columns.findIndex(column => column.toLowerCase() === name.toLowerCase());
        if (index < 0) throw new Error(`Unknown column: ${name}`);
        return row => row[index];
    };
    const condition = (): Predicate => {
        if (take('(')) { const inner = or(); expect(')'); return inner; }
        if (take('NOT')) { const inner = condition(); return row => !inner(row); }
        const left = operand();
        if (take('IS')) {
            const negate = take('NOT'); expect('NULL');
            return row => negate ? left(row) != null : left(row) == null;
        }
        const negate = take('NOT');
        if (take('IN')) {
            expect('('); const values = [literal()];
            while (take(',')) values.push(literal());
            expect(')');
            return row => left(row) != null && (negate !== values.some(value => compare(left(row), value, '=')));
        }
        if (take('LIKE')) {
            const value = literal();
            if (typeof value !== 'string') throw new Error('LIKE requires a text pattern');
            const escaped = value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&').replaceAll('%', '.*').replaceAll('_', '.');
            const regex = new RegExp(`^${escaped}$`, 'su');
            return row => left(row) != null && (negate !== regex.test(String(normalize(left(row)))));
        }
        if (negate) throw new Error('Expected IN or LIKE after NOT');
        const operator = tokens[position++]?.text;
        if (!operator || !['=', '!=', '<>', '>', '<', '>=', '<='].includes(operator)) throw new Error('Expected a comparison, LIKE, IN, or IS NULL');
        const right = operand();
        return row => compare(left(row), right(row), operator);
    };
    const and = (): Predicate => {
        let left = condition();
        while (take('AND')) { const previous = left, right = condition(); left = row => previous(row) && right(row); }
        return left;
    };
    const or = (): Predicate => {
        let left = and();
        while (take('OR')) { const previous = left, right = and(); left = row => previous(row) || right(row); }
        return left;
    };
    if (!tokens.length) return () => true;
    const predicate = or();
    if (position !== tokens.length) throw new Error(`Unexpected token: ${tokens[position]?.text}`);
    return predicate;
}
