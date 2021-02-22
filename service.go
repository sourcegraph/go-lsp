package lsp

import (
	"bytes"
	"encoding/json"
	"strings"
)

type None struct{}

type InitializeParams struct {
	ProcessID int `json:"processId,omitempty"`

	// RootPath is DEPRECATED in favor of the RootURI field.
	RootPath string `json:"rootPath,omitempty"`

	RootURI               DocumentURI        `json:"rootUri,omitempty"`
	ClientInfo            ClientInfo         `json:"clientInfo,omitempty"`
	Trace                 Trace              `json:"trace,omitempty"`
	InitializationOptions interface{}        `json:"initializationOptions,omitempty"`
	Capabilities          ClientCapabilities `json:"capabilities"`

	WorkDoneToken string `json:"workDoneToken,omitempty"`
}

// Root returns the RootURI if set, or otherwise the RootPath with 'file://' prepended.
func (p *InitializeParams) Root() DocumentURI {
	if p.RootURI != "" {
		return p.RootURI
	}
	if strings.HasPrefix(p.RootPath, "file://") {
		return DocumentURI(p.RootPath)
	}
	return DocumentURI("file://" + p.RootPath)
}

type DocumentURI string
type URI string

type ClientInfo struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

type Trace string

type ClientCapabilities struct {
	Workspace    WorkspaceClientCapabilities    `json:"workspace,omitempty"`
	TextDocument TextDocumentClientCapabilities `json:"textDocument,omitempty"`
	Window       WindowClientCapabilities       `json:"window,omitempty"`
	Experimental interface{}                    `json:"experimental,omitempty"`

	// Below are Sourcegraph extensions. They do not live in lspext since
	// they are extending the field InitializeParams.Capabilities

	// XFilesProvider indicates the client provides support for
	// workspace/xfiles. This is a Sourcegraph extension.
	XFilesProvider bool `json:"xfilesProvider,omitempty"`

	// XContentProvider indicates the client provides support for
	// textDocument/xcontent. This is a Sourcegraph extension.
	XContentProvider bool `json:"xcontentProvider,omitempty"`

	// XCacheProvider indicates the client provides support for cache/get
	// and cache/set.
	XCacheProvider bool `json:"xcacheProvider,omitempty"`
}

type WorkspaceClientCapabilities struct {
	WorkspaceEdit struct {
		DocumentChanges    bool     `json:"documentChanges,omitempty"`
		ResourceOperations []string `json:"resourceOperations,omitempty"`
	} `json:"workspaceEdit,omitempty"`

	ApplyEdit bool `json:"applyEdit,omitempty"`

	Symbol struct {
		SymbolKind struct {
			ValueSet []int `json:"valueSet,omitempty"`
		} `json:"symbolKind,omitEmpty"`
	} `json:"symbol,omitempty"`

	ExecuteCommand *struct {
		DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
	} `json:"executeCommand,omitempty"`

	DidChangeWatchedFiles *struct {
		DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
	} `json:"didChangeWatchedFiles,omitempty"`

	WorkspaceFolders bool `json:"workspaceFolders,omitempty"`

	Configuration bool `json:"configuration,omitempty"`
}

type TextDocumentClientCapabilities struct {
	Declaration *struct {
		LinkSupport bool `json:"linkSupport,omitempty"`
	} `json:"declaration,omitempty"`

	Definition *struct {
		LinkSupport bool `json:"linkSupport,omitempty"`
	} `json:"definition,omitempty"`

	Implementation *struct {
		LinkSupport bool `json:"linkSupport,omitempty"`

		DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
	} `json:"implementation,omitempty"`

	TypeDefinition *struct {
		LinkSupport bool `json:"linkSupport,omitempty"`
	} `json:"typeDefinition,omitempty"`

	Synchronization *struct {
		WillSave          bool `json:"willSave,omitempty"`
		DidSave           bool `json:"didSave,omitempty"`
		WillSaveWaitUntil bool `json:"willSaveWaitUntil,omitempty"`
	} `json:"synchronization,omitempty"`

	DocumentSymbol struct {
		SymbolKind struct {
			ValueSet []int `json:"valueSet,omitempty"`
		} `json:"symbolKind,omitEmpty"`

		HierarchicalDocumentSymbolSupport bool `json:"hierarchicalDocumentSymbolSupport,omitempty"`
	} `json:"documentSymbol,omitempty"`

	Formatting *struct {
		DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
	} `json:"formatting,omitempty"`

	RangeFormatting *struct {
		DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
	} `json:"rangeFormatting,omitempty"`

	Rename *struct {
		DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

		PrepareSupport bool `json:"prepareSupport,omitempty"`
	} `json:"rename,omitempty"`

	SemanticHighlightingCapabilities *struct {
		SemanticHighlighting bool `json:"semanticHighlighting,omitempty"`
	} `json:"semanticHighlightingCapabilities,omitempty"`

	CodeAction struct {
		DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

		IsPreferredSupport bool `json:"isPreferredSupport,omitempty"`

		CodeActionLiteralSupport struct {
			CodeActionKind struct {
				ValueSet []CodeActionKind `json:"valueSet,omitempty"`
			} `json:"codeActionKind,omitempty"`
		} `json:"codeActionLiteralSupport,omitempty"`
	} `json:"codeAction,omitempty"`

	Completion struct {
		CompletionItem struct {
			DocumentationFormat []DocumentationFormat `json:"documentationFormat,omitempty"`
			SnippetSupport      bool                  `json:"snippetSupport,omitempty"`
		} `json:"completionItem,omitempty"`

		CompletionItemKind struct {
			ValueSet []CompletionItemKind `json:"valueSet,omitempty"`
		} `json:"completionItemKind,omitempty"`

		ContextSupport bool `json:"contextSupport,omitempty"`
	} `json:"completion,omitempty"`

	SignatureHelp *struct {
		SignatureInformation struct {
			ParameterInformation struct {
				LabelOffsetSupport bool `json:"labelOffsetSupport,omitempty"`
			} `json:"parameterInformation,omitempty"`
		} `json:"signatureInformation,omitempty"`
	} `json:"signatureHelp,omitempty"`

	DocumentLink *struct {
		DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

		TooltipSupport bool `json:"tooltipSupport,omitempty"`
	} `json:"documentLink,omitempty"`

	Hover *struct {
		ContentFormat []string `json:"contentFormat,omitempty"`
	} `json:"hover,omitempty"`

	FoldingRange *struct {
		DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

		RangeLimit interface{} `json:"rangeLimit,omitempty"`

		LineFoldingOnly bool `json:"lineFoldingOnly,omitempty"`
	} `json:"foldingRange,omitempty"`

	CallHierarchy *struct {
		DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
	} `json:"callHierarchy,omitempty"`

	ColorProvider *struct {
		DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
	} `json:"colorProvider,omitempty"`
}

type WindowClientCapabilities struct {
	WorkDoneProgress bool `json:"workDoneProgress,omitempty"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities,omitempty"`
}

type InitializeError struct {
	Retry bool `json:"retry"`
}

type ResourceOperation string

const (
	ROCreate ResourceOperation = "create"
	RODelete ResourceOperation = "delete"
	RORename ResourceOperation = "rename"
)

// TextDocumentSyncKind is a DEPRECATED way to describe how text
// document syncing works. Use TextDocumentSyncOptions instead (or the
// Options field of TextDocumentSyncOptionsOrKind if you need to
// support JSON-(un)marshaling both).
type TextDocumentSyncKind int

const (
	TDSKNone        TextDocumentSyncKind = 0
	TDSKFull        TextDocumentSyncKind = 1
	TDSKIncremental TextDocumentSyncKind = 2
)

type TextDocumentSyncOptions struct {
	OpenClose         bool                 `json:"openClose,omitempty"`
	Change            TextDocumentSyncKind `json:"change"`
	WillSave          bool                 `json:"willSave,omitempty"`
	WillSaveWaitUntil bool                 `json:"willSaveWaitUntil,omitempty"`
	Save              *SaveOptions         `json:"save,omitempty"`
}

// TextDocumentSyncOptions holds either a TextDocumentSyncKind or
// TextDocumentSyncOptions. The LSP API allows either to be specified
// in the (ServerCapabilities).TextDocumentSync field.
type TextDocumentSyncOptionsOrKind struct {
	Kind    *TextDocumentSyncKind
	Options *TextDocumentSyncOptions
}

// MarshalJSON implements json.Marshaler.
func (v *TextDocumentSyncOptionsOrKind) MarshalJSON() ([]byte, error) {
	if v == nil {
		return []byte("null"), nil
	}
	if v.Kind != nil {
		return json.Marshal(v.Kind)
	}
	return json.Marshal(v.Options)
}

// UnmarshalJSON implements json.Unmarshaler.
func (v *TextDocumentSyncOptionsOrKind) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		*v = TextDocumentSyncOptionsOrKind{}
		return nil
	}
	var kind TextDocumentSyncKind
	if err := json.Unmarshal(data, &kind); err == nil {
		// Create equivalent TextDocumentSyncOptions using the same
		// logic as in vscode-languageclient. Also set the Kind field
		// so that JSON-marshaling and unmarshaling are inverse
		// operations (for backward compatibility, preserving the
		// original input but accepting both).
		*v = TextDocumentSyncOptionsOrKind{
			Options: &TextDocumentSyncOptions{OpenClose: true, Change: kind},
			Kind:    &kind,
		}
		return nil
	}
	var tmp TextDocumentSyncOptions
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	*v = TextDocumentSyncOptionsOrKind{Options: &tmp}
	return nil
}

type SaveOptions struct {
	IncludeText bool `json:"includeText"`
}

type WorkspaceFoldersServerCapabilities struct {
	/**
	 * The server has support for workspace folders
	 */
	Supported bool `json:"supported,omitempty"`

	/**
	 * Whether the server wants to receive workspace folder
	 * change notifications.
	 *
	 * If a string is provided, the string is treated as an ID
	 * under which the notification is registered on the client
	 * side. The ID can be used to unregister for these events
	 * using the `client/unregisterCapability` request.
	 */
	ChangeNotifications string `json:"changeNotifications,omitempty"`
}

type FileOperationPatternKind string

const (
	/**
	 * The pattern matches a file only.
	 */
	FOPKFile FileOperationPatternKind = "file"

	/**
	 * The pattern matches a folder only.
	 */
	FOPKFolder FileOperationPatternKind = "folder"
)

type FileOperationPatternOptions struct {

	/**
	 * The pattern should be matched ignoring casing.
	 */
	IgnoreCase bool `json:"IgnoreCase,omitempty"`
}

type FileOperationPattern struct {
	/**
	 * The glob pattern to match. Glob patterns can have the following syntax:
	 * - `*` to match one or more characters in a path segment
	 * - `?` to match on one character in a path segment
	 * - `**` to match any number of path segments, including none
	 * - `{}` to group sub patterns into an OR expression. (e.g. `**​/*.{ts,js}`
	 *   matches all TypeScript and JavaScript files)
	 * - `[]` to declare a range of characters to match in a path segment
	 *   (e.g., `example.[0-9]` to match on `example.0`, `example.1`, …)
	 * - `[!...]` to negate a range of characters to match in a path segment
	 *   (e.g., `example.[!0-9]` to match on `example.a`, `example.b`, but
	 *   not `example.0`)
	 */
	Glob string `json:"glob"`

	/**
	 * Whether to match files or folders with this pattern.
	 *
	 * Matches both if undefined.
	 */
	Matches FileOperationPatternKind `json:"matches"`

	/**
	 * Additional options used during matching.
	 */
	Options FileOperationPatternOptions
}

type FileOperationFilter struct {

	/**
	 * A Uri like `file` or `untitled`.
	 */
	Scheme string `json:"scheme,omitempty"`

	/**
	 * The actual file operation pattern.
	 */
	Pattern FileOperationPattern `json:"pattern,omitempty"`
}

type FileOperationRegistrationOptions struct {
	/**
	 * The actual filters.
	 */
	Filters []FileOperationFilter `json:"filters,omitempty"`
}

type WorkspaceOptionsFileOperations struct {
	/**
	* The server is interested in receiving didCreateFiles
	* notifications.
	 */
	DidCreate *FileOperationRegistrationOptions `json:"didCreate,omitempty"`

	/**
	* The server is interested in receiving willCreateFiles requests.
	 */
	WillCreate *FileOperationRegistrationOptions `json:"willCreate,omitempty"`

	/**
	* The server is interested in receiving didRenameFiles
	* notifications.
	 */
	DidRename *FileOperationRegistrationOptions `json:"didRename,omitempty"`

	/**
	* The server is interested in receiving willRenameFiles requests.
	 */
	WillRename *FileOperationRegistrationOptions `json:"willRename,omitempty"`

	/**
	* The server is interested in receiving didDeleteFiles file
	* notifications.
	 */
	DidDelete *FileOperationRegistrationOptions `json:"didDelete,omitempty"`

	/**
	* The server is interested in receiving willDeleteFiles file
	* requests.
	 */
	WillDelete *FileOperationRegistrationOptions `json:"willDelete,omitempty"`
}

type WorkspaceOptions struct {
	WorkspaceFolders *WorkspaceFoldersServerCapabilities `json:"workspaceFolders,omitempty"`
	/**
	 * The server is interested in file notifications/requests.
	 *
	 * @since 3.16.0
	 */
	FileOperations *WorkspaceOptionsFileOperations `json:"fileOperations,omitempty"`
}

type ServerCapabilities struct {
	TextDocumentSync                 *TextDocumentSyncOptionsOrKind   `json:"textDocumentSync,omitempty"`
	CompletionProvider               *CompletionOptions               `json:"completionProvider,omitempty"`
	HoverProvider                    bool                             `json:"hoverProvider,omitempty"`
	SignatureHelpProvider            *SignatureHelpOptions            `json:"signatureHelpProvider,omitempty"`
	DeclarationProvider              *DeclarationOptions              `json:"declarationProvider,omitempty"`
	DefinitionProvider               bool                             `json:"definitionProvider,omitempty"`
	TypeDefinitionProvider           bool                             `json:"typeDefinitionProvider,omitempty"`
	ImplementationProvider           *ImplementationOptions           `json:"implementationProvider,omitempty"`
	ReferencesProvider               bool                             `json:"referencesProvider,omitempty"`
	DocumentHighlightProvider        bool                             `json:"documentHighlightProvider,omitempty"`
	DocumentSymbolProvider           bool                             `json:"documentSymbolProvider,omitempty"`
	CodeActionProvider               bool                             `json:"codeActionProvider,omitempty"`
	CodeLensProvider                 *CodeLensOptions                 `json:"codeLensProvider,omitempty"`
	DocumentLinkProvider             *DocumentLinkOptions             `json:"documentLinkProvider,omitempty"`
	ColorProvider                    *DocumentColorOptions            `json:"colorProvider,omitempty"`
	DocumentFormattingProvider       bool                             `json:"documentFormattingProvider,omitempty"`
	DocumentRangeFormattingProvider  bool                             `json:"documentRangeFormattingProvider,omitempty"`
	DocumentOnTypeFormattingProvider *DocumentOnTypeFormattingOptions `json:"documentOnTypeFormattingProvider,omitempty"`
	RenameProvider                   bool                             `json:"renameProvider,omitempty"`
	FoldingRangeProvider             *FoldingRangeOptions             `json:"foldingRangeProvider,omitempty"`
	ExecuteCommandProvider           *ExecuteCommandOptions           `json:"executeCommandProvider,omitempty"`
	SelectionRangeProvider           *SelectionRangeOptions           `json:"selectionRangeProvider,omitempty"`
	LinkedEditingRangeProvider       *LinkedEditingRangeOptions       `json:"linkedEditingRangeProvider,omitempty"`
	CallHierarchyProvider            *CallHierarchyOptions            `json:"callHierarchyProvider,omitempty"`
	SemanticTokensProvider           *SemanticTokensOptions           `json:"semanticTokensProvider,omitempty"`
	MonikerProvider                  *MonikerOptions                  `json:"monikerProvider,omitempty"`
	WorkspaceSymbolProvider          bool                             `json:"workspaceSymbolProvider,omitempty"`
	Workspace                        *WorkspaceOptions                `json:"workspace,omitempty"`

	// XWorkspaceReferencesProvider indicates the server provides support for
	// xworkspace/references. This is a Sourcegraph extension.
	XWorkspaceReferencesProvider bool `json:"xworkspaceReferencesProvider,omitempty"`

	// XDefinitionProvider indicates the server provides support for
	// textDocument/xdefinition. This is a Sourcegraph extension.
	XDefinitionProvider bool `json:"xdefinitionProvider,omitempty"`

	// XWorkspaceSymbolByProperties indicates the server provides support for
	// querying symbols by properties with WorkspaceSymbolParams.symbol. This
	// is a Sourcegraph extension.
	XWorkspaceSymbolByProperties bool `json:"xworkspaceSymbolByProperties,omitempty"`

	Experimental interface{} `json:"experimental,omitempty"`
}

type CompletionOptions struct {
	ResolveProvider   bool     `json:"resolveProvider,omitempty"`
	TriggerCharacters []string `json:"triggerCharacters,omitempty"`
}

type DocumentOnTypeFormattingOptions struct {
	FirstTriggerCharacter string   `json:"firstTriggerCharacter"`
	MoreTriggerCharacter  []string `json:"moreTriggerCharacter,omitempty"`
}

type CodeLensOptions struct {
	ResolveProvider bool `json:"resolveProvider,omitempty"`
}

type WorkDoneProgressOptions struct {
	WorkDoneProgress bool `json:"workDoneProgress,omitempty"`
}

type FoldingRangeOptions struct {
	WorkDoneProgressOptions
	ResolveProvider bool `json:"resolveProvider,omitempty"`
}

type SelectionRangeOptions struct {
	WorkDoneProgressOptions
}

type LinkedEditingRangeOptions struct {
	WorkDoneProgressOptions
}

type CallHierarchyOptions struct {
	WorkDoneProgressOptions
}

type SemanticTokenOptionsFull struct {
	Delta bool `json:"delta,omitempty"`
}

type SemanticTokensLegend struct {
	/**
	 * The token types a server uses.
	 */
	TokenTypes []string `json:"tokenTypes,omitempty"`

	/**
	 * The token modifiers a server uses.
	 */
	TokenModifiers []string `json:"tokenModifiers,omitempty"`
}

type SemanticTokensOptions struct {
	WorkDoneProgressOptions

	/**
	 * The legend used by the server
	 */
	Legend SemanticTokensLegend `json:"legend,omitempty"`

	/**
	 * Server supports providing semantic tokens for a specific range
	 * of a document.
	 */
	Range bool `json:"range,omitempty"`

	/**
	 * Server supports providing semantic tokens for a full document.
	 */
	Full *SemanticTokenOptionsFull `json:"full,omitempty"`
}

type MonikerOptions struct {
	WorkDoneProgressOptions
}

type DocumentLinkOptions struct {
	ResolveProvider bool `json:"resolveProvider,omitempty"`
}

type DocumentColorOptions struct {
	ResolveProvider bool `json:"resolveProvider,omitempty"`
}

type SignatureHelpOptions struct {
	TriggerCharacters []string `json:"triggerCharacters,omitempty"`
}

type ExecuteCommandOptions struct {
	Commands []string `json:"commands"`
}

type ExecuteCommandParams struct {
	Command   string        `json:"command"`
	Arguments []interface{} `json:"arguments,omitempty"`
}

type DeclarationOptions struct {
}

type ImplementationOptions struct {
}

type CompletionItemKind int

const (
	_ CompletionItemKind = iota
	CIKText
	CIKMethod
	CIKFunction
	CIKConstructor
	CIKField
	CIKVariable
	CIKClass
	CIKInterface
	CIKModule
	CIKProperty
	CIKUnit
	CIKValue
	CIKEnum
	CIKKeyword
	CIKSnippet
	CIKColor
	CIKFile
	CIKReference
	CIKFolder
	CIKEnumMember
	CIKConstant
	CIKStruct
	CIKEvent
	CIKOperator
	CIKTypeParameter
)

func (c CompletionItemKind) String() string {
	return completionItemKindName[c]
}

var completionItemKindName = map[CompletionItemKind]string{
	CIKText:          "text",
	CIKMethod:        "method",
	CIKFunction:      "function",
	CIKConstructor:   "constructor",
	CIKField:         "field",
	CIKVariable:      "variable",
	CIKClass:         "class",
	CIKInterface:     "interface",
	CIKModule:        "module",
	CIKProperty:      "property",
	CIKUnit:          "unit",
	CIKValue:         "value",
	CIKEnum:          "enum",
	CIKKeyword:       "keyword",
	CIKSnippet:       "snippet",
	CIKColor:         "color",
	CIKFile:          "file",
	CIKReference:     "reference",
	CIKFolder:        "folder",
	CIKEnumMember:    "enumMember",
	CIKConstant:      "constant",
	CIKStruct:        "struct",
	CIKEvent:         "event",
	CIKOperator:      "operator",
	CIKTypeParameter: "typeParameter",
}

type CompletionItemTag int

const (
	CITDeprecated CompletionItemTag = 1
)

type InsertTextMode int

const (
	/**
	 * The insertion or replace strings is taken as it is. If the
	 * value is multi line the lines below the cursor will be
	 * inserted using the indentation defined in the string value.
	 * The client will not apply any kind of adjustments to the
	 * string.
	 */
	ITMAsIs InsertTextMode = 1

	/**
	 * The editor adjusts leading whitespace of new lines so that
	 * they match the indentation up to the cursor of the line for
	 * which the item is accepted.
	 *
	 * Consider a line like this: <2tabs><cursor><3tabs>foo. Accepting a
	 * multi line completion item is indented using 2 tabs and all
	 * following lines inserted will be indented using 2 tabs as well.
	 */
	ITMAdjustIndentation = 2
)

type CompletionItem struct {
	Label               string              `json:"label"`
	Kind                CompletionItemKind  `json:"kind,omitempty"`
	Tags                []CompletionItemTag `json:"tags,omitempty"`
	Detail              string              `json:"detail,omitempty"`
	Documentation       string              `json:"documentation,omitempty"`
	Preselect           bool                `json:"preselect,omitempty"`
	SortText            string              `json:"sortText,omitempty"`
	FilterText          string              `json:"filterText,omitempty"`
	InsertText          string              `json:"insertText,omitempty"`
	InsertTextFormat    InsertTextFormat    `json:"insertTextFormat,omitempty"`
	InsertTextMode      *InsertTextMode     `json:"insertTextMode,omitempty"`
	TextEdit            *TextEdit           `json:"textEdit,omitempty"`
	AdditionalTextEdits []*TextEdit         `json:"additionalTextEdits,omitempty"`
	CommitCharacters    []string            `json:"commitCharacters,omitempty"`
	Command             *Command            `json:"command,omitempty"`
	Data                interface{}         `json:"data,omitempty"`
}

type CompletionList struct {
	IsIncomplete bool             `json:"isIncomplete"`
	Items        []CompletionItem `json:"items"`
}

type CompletionTriggerKind int

const (
	CTKInvoked          CompletionTriggerKind = 1
	CTKTriggerCharacter                       = 2
)

type DocumentationFormat string

const (
	DFPlainText DocumentationFormat = "plaintext"
)

type InsertTextFormat int

const (
	ITFPlainText InsertTextFormat = 1
	ITFSnippet                    = 2
)

type CompletionContext struct {
	TriggerKind      CompletionTriggerKind `json:"triggerKind"`
	TriggerCharacter string                `json:"triggerCharacter,omitempty"`
}

type CompletionParams struct {
	TextDocumentPositionParams
	Context CompletionContext `json:"context,omitempty"`
}

type DocumentHighlightParams struct {
	TextDocumentPositionParams
	Context CompletionContext `json:"context,omitempty"`
}

type Hover struct {
	Contents MarkupContent `json:"contents"`
	Range    *Range        `json:"range,omitempty"`
}

type hover Hover

func (h Hover) MarshalJSON() ([]byte, error) {
	return json.Marshal(hover(h))
}

type MarkupKind string

const (
	MUKPlainText = "plaintext"
	MUKMarkdown  = "markdown"
)

type MarkupContent markupContent
type markupContent struct {
	Kind  MarkupKind `json:"kind"`
	Value string     `json:"value"`
}

type SignatureHelp struct {
	Signatures      []SignatureInformation `json:"signatures"`
	ActiveSignature int                    `json:"activeSignature"`
	ActiveParameter int                    `json:"activeParameter"`
}

type SignatureInformation struct {
	Label         string                 `json:"label"`
	Documentation string                 `json:"documentation,omitempty"`
	Parameters    []ParameterInformation `json:"parameters,omitempty"`
}

type ParameterInformation struct {
	Label         string `json:"label"`
	Documentation string `json:"documentation,omitempty"`
}

type ReferenceContext struct {
	IncludeDeclaration bool `json:"includeDeclaration"`

	// Sourcegraph extension
	XLimit int `json:"xlimit,omitempty"`
}

type ReferenceParams struct {
	TextDocumentPositionParams
	Context ReferenceContext `json:"context"`
}

type DocumentHighlightKind int

const (
	Text  DocumentHighlightKind = 1
	Read                        = 2
	Write                       = 3
)

type DocumentHighlight struct {
	Range Range `json:"range"`
	Kind  int   `json:"kind,omitempty"`
}

type DocumentSymbolParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type SymbolKind int

// The SymbolKind values are defined at https://microsoft.github.io/language-server-protocol/specification.
const (
	SKFile          SymbolKind = 1
	SKModule        SymbolKind = 2
	SKNamespace     SymbolKind = 3
	SKPackage       SymbolKind = 4
	SKClass         SymbolKind = 5
	SKMethod        SymbolKind = 6
	SKProperty      SymbolKind = 7
	SKField         SymbolKind = 8
	SKConstructor   SymbolKind = 9
	SKEnum          SymbolKind = 10
	SKInterface     SymbolKind = 11
	SKFunction      SymbolKind = 12
	SKVariable      SymbolKind = 13
	SKConstant      SymbolKind = 14
	SKString        SymbolKind = 15
	SKNumber        SymbolKind = 16
	SKBoolean       SymbolKind = 17
	SKArray         SymbolKind = 18
	SKObject        SymbolKind = 19
	SKKey           SymbolKind = 20
	SKNull          SymbolKind = 21
	SKEnumMember    SymbolKind = 22
	SKStruct        SymbolKind = 23
	SKEvent         SymbolKind = 24
	SKOperator      SymbolKind = 25
	SKTypeParameter SymbolKind = 26
)

func (s SymbolKind) String() string {
	return symbolKindName[s]
}

var symbolKindName = map[SymbolKind]string{
	SKFile:          "File",
	SKModule:        "Module",
	SKNamespace:     "Namespace",
	SKPackage:       "Package",
	SKClass:         "Class",
	SKMethod:        "Method",
	SKProperty:      "Property",
	SKField:         "Field",
	SKConstructor:   "Constructor",
	SKEnum:          "Enum",
	SKInterface:     "Interface",
	SKFunction:      "Function",
	SKVariable:      "Variable",
	SKConstant:      "Constant",
	SKString:        "String",
	SKNumber:        "Number",
	SKBoolean:       "Boolean",
	SKArray:         "Array",
	SKObject:        "Object",
	SKKey:           "Key",
	SKNull:          "Null",
	SKEnumMember:    "EnumMember",
	SKStruct:        "Struct",
	SKEvent:         "Event",
	SKOperator:      "Operator",
	SKTypeParameter: "TypeParameter",
}

type SymbolTag int

const (
	SYTDeprecated SymbolTag = 1
)

type SymbolInformation struct {
	Name          string     `json:"name"`
	Kind          SymbolKind `json:"kind"`
	Location      Location   `json:"location"`
	ContainerName string     `json:"containerName,omitempty"`
}

type DocumentSymbol struct {
	Name           string           `json:"name"`
	Detail         string           `json:"detail,omitempty"` // Usually the signature of the function
	Kind           SymbolKind       `json:"kind"`
	Tags           []SymbolTag      `json:"tags,omitempty"`
	Range          Range            `json:"range"`
	SelectionRange Range            `json:"selectionRange"`
	Children       []DocumentSymbol `json:"children,omitempty"` // For example fields in a class
}

type WorkspaceSymbolParams struct {
	Query string `json:"query"`
	//Limit int    `json:"limit"`
	//WorkDoneToken
	//PartialResultToken
}

type ConfigurationParams struct {
	Items []ConfigurationItem `json:"items"`
}

type ConfigurationItem struct {
	ScopeURI string `json:"scopeUri,omitempty"`
	Section  string `json:"section,omitempty"`
}

type ConfigurationResult []interface{}

type CodeActionContext struct {
	Diagnostics []Diagnostic `json:"diagnostics"`
}

type CodeActionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Range        Range                  `json:"range"`
	Context      CodeActionContext      `json:"context"`
}

type CodeLensParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type CodeLens struct {
	Range   Range       `json:"range"`
	Command Command     `json:"command,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type DocumentFormattingParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Options      FormattingOptions      `json:"options"`
}

type FormattingOptions struct {
	TabSize      int    `json:"tabSize"`
	InsertSpaces bool   `json:"insertSpaces"`
	Key          string `json:"key"`
}

type RenameParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
	NewName      string                 `json:"newName"`
}

type DidOpenTextDocumentParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}

type DidChangeTextDocumentParams struct {
	TextDocument   VersionedTextDocumentIdentifier  `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

type TextDocumentContentChangeEvent struct {
	Range       *Range `json:"range,omitEmpty"`
	RangeLength uint   `json:"rangeLength,omitEmpty"`
	Text        string `json:"text"`
}

type DidCloseTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type DidSaveTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type MessageType int

const (
	MTError   MessageType = 1
	MTWarning             = 2
	MTInfo                = 3
	MTLog                 = 4
)

type ShowMessageParams struct {
	Type    MessageType `json:"type"`
	Message string      `json:"message"`
}

type MessageActionItem struct {
	Title string `json:"title"`
}

type ShowMessageRequestParams struct {
	Type    MessageType         `json:"type"`
	Message string              `json:"message"`
	Actions []MessageActionItem `json:"actions"`
}

type LogMessageParams struct {
	Type    MessageType `json:"type"`
	Message string      `json:"message"`
}

type DidChangeConfigurationParams struct {
	Settings interface{} `json:"settings"`
}

type FileChangeType int

const (
	FCTCreated FileChangeType = 1
	FCTChanged                = 2
	FCTDeleted                = 3
)

type FileEvent struct {
	URI  DocumentURI `json:"uri"`
	Type int         `json:"type"`
}

type DidChangeWatchedFilesParams struct {
	Changes []FileEvent `json:"changes"`
}

type PublishDiagnosticsParams struct {
	URI         DocumentURI  `json:"uri"`
	Version     uint         `json:"version"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

type DocumentRangeFormattingParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Range        Range                  `json:"range"`
	Options      FormattingOptions      `json:"options"`
}

type DocumentOnTypeFormattingParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
	Ch           string                 `json:"ch"`
	Options      FormattingOptions      `json:"formattingOptions"`
}

type CancelParams struct {
	ID ID `json:"id"`
}

/*
type SemanticHighlightingParams struct {
	TextDocument VersionedTextDocumentIdentifier   `json:"textDocument"`
	Lines        []SemanticHighlightingInformation `json:"lines"`
}

// SemanticHighlightingInformation represents a semantic highlighting
// information that has to be applied on a specific line of the text
// document.
type SemanticHighlightingInformation struct {
	// Line is the zero-based line position in the text document.
	Line int `json:"line"`

	// Tokens is a base64 encoded string representing every single highlighted
	// characters with its start position, length and the "lookup table" index of
	// the semantic highlighting [TextMate scopes](https://manual.macromates.com/en/language_grammars).
	// If the `tokens` is empty or not defined, then no highlighted positions are
	// available for the line.
	Tokens SemanticHighlightingTokens `json:"tokens,omitempty"`
}

type semanticHighlightingInformation struct {
	Line   int     `json:"line"`
	Tokens *string `json:"tokens"`
}

// MarshalJSON implements json.Marshaler.
func (v *SemanticHighlightingInformation) MarshalJSON() ([]byte, error) {
	tokens := string(v.Tokens.Serialize())
	return json.Marshal(&semanticHighlightingInformation{
		Line:   v.Line,
		Tokens: &tokens,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (v *SemanticHighlightingInformation) UnmarshalJSON(data []byte) error {
	var info semanticHighlightingInformation
	err := json.Unmarshal(data, &info)
	if err != nil {
		return err
	}

	if info.Tokens != nil {
		v.Tokens, err = DeserializeSemanticHighlightingTokens([]byte(*info.Tokens))
		if err != nil {
			return err
		}
	}

	v.Line = info.Line
	return nil
}

type SemanticHighlightingTokens []SemanticHighlightingToken

func (v SemanticHighlightingTokens) Serialize() []byte {
	var chunks [][]byte

	// Writes each token to `tokens` in the byte format specified by the LSP
	// proposal. Described below:
	// |<---- 4 bytes ---->|<-- 2 bytes -->|<--- 2 bytes -->|
	// |    character      |  length       |    index       |
	for _, token := range v {
		chunk := make([]byte, 8)
		binary.BigEndian.PutUint32(chunk[:4], token.Character)
		binary.BigEndian.PutUint16(chunk[4:6], token.Length)
		binary.BigEndian.PutUint16(chunk[6:], token.Scope)
		chunks = append(chunks, chunk)
	}

	src := make([]byte, len(chunks)*8)
	for i, chunk := range chunks {
		copy(src[i*8:i*8+8], chunk)
	}

	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	return dst
}

func DeserializeSemanticHighlightingTokens(src []byte) (SemanticHighlightingTokens, error) {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	n, err := base64.StdEncoding.Decode(dst, src)
	if err != nil {
		return nil, err
	}

	var chunks [][]byte
	for i := 7; i < len(dst[:n]); i += 8 {
		chunks = append(chunks, dst[i-7:i+1])
	}

	var tokens SemanticHighlightingTokens
	for _, chunk := range chunks {
		tokens = append(tokens, SemanticHighlightingToken{
			Character: binary.BigEndian.Uint32(chunk[:4]),
			Length:    binary.BigEndian.Uint16(chunk[4:6]),
			Scope:     binary.BigEndian.Uint16(chunk[6:]),
		})
	}

	return tokens, nil
}

type SemanticHighlightingToken struct {
	Character uint32
	Length    uint16
	Scope     uint16
}
*/
