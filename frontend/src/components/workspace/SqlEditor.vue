<template>
    <div ref="editorContainer" class="h-full w-full overflow-hidden" />
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch } from "vue";

import { EditorView, keymap } from "@codemirror/view";
import { EditorState } from "@codemirror/state";
import { defaultKeymap, indentWithTab } from "@codemirror/commands";

import { sql, PostgreSQL } from "@codemirror/lang-sql";

import { oneDark } from "@codemirror/theme-one-dark";

const props = defineProps<{
    modelValue: string;
}>();

const emit = defineEmits<{
    "update:modelValue": [value: string];
    execute: [];
}>();

const editorContainer = ref<HTMLElement | null>(null);

let editorView: EditorView | null = null;

onMounted(() => {
    if (!editorContainer.value) {
        return;
    }

    const startState = EditorState.create({
        doc: props.modelValue,

        extensions: [
            // SQL syntax highlighting
            sql({
                dialect: PostgreSQL,
            }),

            // Dark editor theme
            oneDark,

            // Basic editing
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

            // Ctrl/Cmd + Enter → execute
            keymap.of([
                {
                    key: "Mod-Enter",

                    run: () => {
                        emit("execute");
                        return true;
                    },
                },
            ]),
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

onBeforeUnmount(() => {
    editorView?.destroy();
    editorView = null;
});
</script>
