export type IndexedResultRow = { row: unknown[]; index: number };

function sortableValue(value: unknown): unknown {
    if (typeof value !== 'string') return value;
    try {
        const parsed = JSON.parse(value);
        if (parsed && typeof parsed === 'object') {
            if ('$numberLong' in parsed) return BigInt(parsed.$numberLong);
            if ('$numberInt' in parsed) return Number(parsed.$numberInt);
            if ('$numberDouble' in parsed) return Number(parsed.$numberDouble);
            if ('$oid' in parsed) return parsed.$oid;
        }
    } catch { /* Ordinary text keeps its original value. */ }
    return value;
}

// Preserve original row indexes for the grid's row keys and row numbers.
export function sortResultRows(rows: IndexedResultRow[], column: number, direction?: 'asc' | 'desc'): IndexedResultRow[] {
    if (column < 1 || !direction) return rows;
    const multiplier = direction === 'desc' ? -1 : 1;
    return [...rows].sort((a, b) => {
        const left = sortableValue(a.row[column - 1]);
        const right = sortableValue(b.row[column - 1]);
        let comparison = 0;
        if (left == null || right == null) comparison = left == null ? (right == null ? 0 : -1) : 1;
        else if ((typeof left === 'number' || typeof left === 'bigint') && (typeof right === 'number' || typeof right === 'bigint')) {
            comparison = left < right ? -1 : left > right ? 1 : 0;
        } else {
            const text = (value: unknown) => typeof value === 'object' ? JSON.stringify(value) : String(value);
            comparison = text(left).localeCompare(text(right));
        }
        return comparison * multiplier || a.index - b.index;
    });
}
