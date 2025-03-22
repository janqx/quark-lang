package quark

type SymbolScope uint8

const (
	ScopeLocal SymbolScope = iota
	ScopeOuter
	ScopeGlobal
)

type SymbolTableType uint8

const (
	TypeGlobal SymbolTableType = iota
	TypeBlock
	TypeFunction
)

type Symbol struct {
	Name       string
	Index      int
	Scope      SymbolScope
	OuterIndex int
	OuterScope SymbolScope
	Owner      *SymbolTable
}

type SymbolTable struct {
	Parent      *SymbolTable
	Owner       *SymbolTable
	Type        SymbolTableType
	Symbols     map[string]*Symbol
	LocalCount  int
	OuterCount  int
	GlobalCount int
}

func NewSymbolTable(parent *SymbolTable, stt SymbolTableType) *SymbolTable {
	return &SymbolTable{
		Parent:      parent,
		Owner:       nil,
		Type:        stt,
		Symbols:     make(map[string]*Symbol),
		LocalCount:  0,
		OuterCount:  0,
		GlobalCount: 0,
	}
}

func (s *SymbolTable) FindSymbol(name string) *Symbol {
	return s.Symbols[name]
}

func (s *SymbolTable) AddLocalSymbol(name string) *Symbol {
	if symbol, ok := s.Symbols[name]; ok && symbol.Scope == ScopeLocal {
		return symbol
	}
	symbol := &Symbol{
		Name:  name,
		Index: s.LocalCount,
		Scope: ScopeLocal,
		Owner: s,
	}
	s.Symbols[name] = symbol
	s.LocalCount++
	return symbol
}

func (s *SymbolTable) AddGlobalSymbol(name string) *Symbol {
	if s.Parent != nil {
		panic("must be global symbol table")
	}
	if symbol, ok := s.Symbols[name]; ok && symbol.Scope == ScopeGlobal {
		return symbol
	}
	symbol := &Symbol{
		Name:  name,
		Index: s.GlobalCount,
		Scope: ScopeGlobal,
		Owner: s,
	}
	s.Symbols[name] = symbol
	s.GlobalCount++
	return symbol
}

func (s *SymbolTable) Push(stt SymbolTableType) *SymbolTable {
	if stt == TypeGlobal {
		panic("can't push a global symbol table")
	}

	st := NewSymbolTable(s, stt)

	if stt != TypeFunction {
		st.LocalCount = s.LocalCount
		st.OuterCount = s.OuterCount
	}

	if stt == TypeBlock {
		if s.Type == TypeFunction {
			st.Owner = s
		} else {
			st.Owner = s.Owner
		}
	}

	for name, symbol := range s.Symbols {
		if stt == TypeBlock || symbol.Scope == ScopeGlobal {
			st.Symbols[name] = symbol
		} else if stt == TypeFunction {
			st.Symbols[name] = &Symbol{
				Name:       name,
				Index:      st.OuterCount,
				Scope:      ScopeOuter,
				OuterIndex: symbol.Index,
				OuterScope: symbol.Scope,
			}
			st.OuterCount++
		}
	}

	return st
}

func (s *SymbolTable) Pop() *SymbolTable {
	parent := s.Parent
	if parent == nil {
		return nil
	}

	if s.Type == TypeBlock {
		parent.LocalCount = s.LocalCount
	}

	return parent
}
