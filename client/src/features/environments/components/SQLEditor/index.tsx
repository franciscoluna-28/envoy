import CodeMirror from '@uiw/react-codemirror';
import { sql } from '@codemirror/lang-sql';

type Props = {
    value: string;
    onChange: (value: string) => void;
}

export function SQLEditor({ value, onChange }: Props) {
  return (
    <div className="border border-stone-200 rounded-none h-full w-full flex flex-col overflow-hidden bg-white">
      <CodeMirror
        value={value}
        height="100%" 
        minHeight="100%"
        theme="light" 
        extensions={[sql()]}
        onChange={onChange}
        basicSetup={{
          lineNumbers: true,
          foldGutter: false,
          dropCursor: true,
          allowMultipleSelections: false,
          indentOnInput: true,
          highlightActiveLine: true,
        }}
        className="text-sm font-mono flex-1 overflow-auto custom-sql-editor"
      />

      <style dangerouslySetInnerHTML={{ __html: `
        .custom-sql-editor .cm-editor {
          height: 100% !important;
        }
        .custom-sql-editor .cm-scroller {
          overflow: auto !important;
        }
      `}} />
    </div>
  );
}