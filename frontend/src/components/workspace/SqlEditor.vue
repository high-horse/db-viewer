<template>
     <div ref="editorContainer" class="h-full w-full overflow-hidden" /> 
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch } from "vue";
import { Prec, Compartment } from "@codemirror/state";
import { linter, type Diagnostic } from "@codemirror/lint";
import { syntaxTree } from "@codemirror/language";
import { EditorView, keymap } from "@codemirror/view";
import { EditorState } from "@codemirror/state";
import { defaultKeymap, indentWithTab } from "@codemirror/commands";
import { EditorSelection } from "@codemirror/state";
import { sql, PostgreSQL, MySQL, SQLite } from "@codemirror/lang-sql";
import { nord } from "@fsegurai/codemirror-theme-nord";
import { materialDark } from "@fsegurai/codemirror-theme-material-dark";
import { oneDark } from "@codemirror/theme-one-dark";
import { autocompletion } from "@codemirror/autocomplete";

const props = defineProps<{
    modelValue: string;
    dbDriver: "pgx" | "mysql" | "sqlite";
}>();

const emit = defineEmits<{
    "update:modelValue": [value: string];
    execute: [sql: string];
}>();

const editorContainer = ref<HTMLElement | null>(null);

let editorView: EditorView | null = null;
const sqlDialect = new Compartment();

function getDialect(driver: string) {
    switch (driver) {
        case "pgx":
            return PostgreSQL;
        case "mysql":
            return MySQL;
        case "sqlite":
            return SQLite;
        default:
            return PostgreSQL;
    }
}

function sqlSyntaxLinter(view: EditorView): Diagnostic[] {
    const diagnostics: Diagnostic[] = [];

    syntaxTree(view.state).iterate({
        enter(node) {
            if (node.type.isError) {
                diagnostics.push({
                    from: node.from,
                    to: Math.max(node.to, node.from + 1),
                    severity: "error",
                    message: "SQL syntax error",
                });
            }
        },
    });

    return diagnostics;
}

function getCurrentQuery(view: EditorView): string | null {
    const { from, to, head } = view.state.selection.main;

    // 1. If user has selected text, execute the selection
    if (from !== to) {
        const selected = view.state.doc.sliceString(from, to).trim();
        return selected || null;
    }

    // 2. Use the syntax tree to find the Statement at cursor
    const tree = syntaxTree(view.state);
    const statements = tree.topNode.getChildren("Statement");

    for (const stmt of statements) {
        if (head >= stmt.from && head <= stmt.to) {
            view.dispatch({
                selection: EditorSelection.range(stmt.from, stmt.to),
            });
            return (
                view.state.doc.sliceString(stmt.from, stmt.to).trim() || null
            );
        }
    }

    // 3. Cursor is between statements (e.g. blank line). Optional: pick nearest.
    // Uncomment below if you want "run nearest" instead of "do nothing".
    /*
  let nearest: (typeof statements)[0] | null = null;
  let minDist = Infinity;
  for (const stmt of statements) {
    const dist = Math.min(Math.abs(head - stmt.from), Math.abs(head - stmt.to));
    if (dist < minDist) {
      minDist = dist;
      nearest = stmt;
    }
  }
  if (nearest && minDist <= 2) {
    return view.state.doc.sliceString(nearest.from, nearest.to).trim() || null;
  }
  */

    return null;
}

onMounted(() => {
    if (!editorContainer.value) {
        return;
    }

    const startState = EditorState.create({
        doc: props.modelValue,

        extensions: [
            // SQL syntax highlighting
            sqlDialect.of(
                sql({
                    dialect: getDialect(props.dbDriver),
                    upperCaseKeywords: true,
                }),
            ),

            autocompletion({
                activateOnTyping: false, // disable automatic completion to reduce CPU
                defaultKeymap: true, // Ctrl+Space to force-open
                closeOnBlur: true,
            }),
            // linter(sqlSyntaxLinter, {
            //     delay: 300,
            // }),
            // 
            // Dark editor theme
            materialDark, // oneDark, // nord,
            // Basic editing
            //
            keymap.of([...defaultKeymap, indentWithTab]),

            // Custom styling
            EditorView.theme({
                "&": {
                    height: "100%",
                    backgroundColor: "#0c0b09",
                    color: "#d1d5db",
                },

                ".cm-content": {
                    padding: "16px",
                    fontFamily:
                        "ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace",
                    fontSize: "12px",
                    lineHeight: "1.7",
                    caretColor: "#f59e0b",
                },

                ".cm-scroller": {
                    overflow: "auto",
                },

                ".cm-gutters": {
                    backgroundColor: "#100e0c",
                    color: "#4b4540",
                    border: "none",
                },

                ".cm-activeLineGutter": {
                    backgroundColor: "#231f1a",
                    color: "#f59e0b",
                },

                ".cm-activeLine": {
                    backgroundColor: "rgba(245, 158, 11, 0.04)",
                },

                ".cm-selectionBackground": {
                    backgroundColor: "rgba(245, 158, 11, 0.20) !important",
                },

                ".cm-cursor": {
                    borderLeftColor: "#f59e0b",
                },
            }),

            // Emit changes back to Pinia
            EditorView.updateListener.of((update) => {
                if (!update.docChanged) {
                    return;
                }

                emit("update:modelValue", update.state.doc.toString());
            }),
            Prec.highest(
                keymap.of([
                    {
                        key: "Mod-Enter",

                        run: () => {
                            if (!editorView) return true;

                            const sql = getCurrentQuery(editorView);
                            if (!sql) {
                                return true;
                            }

                            console.log("Executing SQL:", sql);
                            emit("execute", sql);
                            return true;
                        },
                    },
                ]),
            ),
        ],
    });

    editorView = new EditorView({
        state: startState,
        parent: editorContainer.value,
    });
});

watch(
    () => props.modelValue,
    (newValue) => {
        if (!editorView) {
            return;
        }

        const currentValue = editorView.state.doc.toString();

        if (currentValue === newValue) {
            return;
        }

        editorView.dispatch({
            changes: {
                from: 0,
                to: editorView.state.doc.length,
                insert: newValue,
            },
        });
    },
);

watch(
    () => props.dbDriver,
    (driver) => {
        if (!editorView) {
            return;
        }

        editorView.dispatch({
            effects: sqlDialect.reconfigure(
                sql({
                    dialect: getDialect(driver),
                }),
            ),
        });
    },
);

onBeforeUnmount(() => {
    editorView?.destroy();
    editorView = null;
});
</script>
