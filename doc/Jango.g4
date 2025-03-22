// java -cp antlr-4.13.1-complete.jar org.antlr.v4.Tool -Dlanguage=Go -o ../antlr -package antlr -visitor -no-listener Jango.g4
grammar Jango;

@parser::members {
func (p *JangoParser) lineTerminatorAhead() bool {
	possibleIndexEosToken := p.GetCurrentToken().GetTokenIndex() - 1
	ahead := p.GetTokenStream().Get(possibleIndexEosToken)
	if ahead.GetChannel() != antlr.LexerHidden {
		return false
	}
  if ahead.GetTokenType() == JangoParserWhiteSpace {
    return true
  }
	return ahead.GetTokenType() == JangoParserWhiteSpace
}

func (p *JangoParser) closeBrace() bool {
	return p.GetTokenStream().LT(1).GetTokenType() == JangoParserCloseBrace
}
}

fragment DIGIT: [0-9];

MultiLineComment:   '/*' .*? '*/'                  -> channel(HIDDEN);
SingleLineComment:  '//' ~[\r\n\u2028\u2029]*      -> channel(HIDDEN);
WhiteSpace:         [\t\u000B\u000C\u0020\u00A0]+  -> channel(HIDDEN);
LineTerminator:     [\r\n\u2028\u2029]             -> channel(HIDDEN);

Identifier: [_a-zA-Z][_a-zA-Z0-9]*;

FloatLiteral: DIGIT+ '.' DIGIT+;
IntLiteral: DIGIT+;
StringLiteral
  : '"' (~[\r\n\u2028\u2029]*?) '"'
  | '\'' (~[\r\n\u2028\u2029]*?) '\''
  | '`' (.*?) '`'
  ;

OpenBracket: '[';
CloseBracket: ']';
OpenParen: '(';
CloseParen: ')';
OpenBrace: '{';
CloseBrace: '}';
Comma: ',';
SemiColon: ';';
Colon: ':';
Assign: '=';
Dot: '.';
Plus: '+';
Minus: '-';
BitNot: '~';
BitAnd: '&';
BitOr: '|';
BitXor: '^';
Not: '!';
Multiply: '*';
Divide: '/';
Modulus: '%';
LogicAnd: '&&';
LogicOr: '||';
LessThan: '<';
LessThanEquals: '<=';
GreaterThan: '>';
GreaterThanEquals: '>=';
Equals: '==';
NotEquals: '!=';

file
  : statementList? EOF
  ;

statementList
  : statement+
  ;

statement
  : importStatement eos__
  | assignStatement eos__
  | returnStatement eos__
  | callFunctionStatement eos__
  | ifStatement
  | forStatement
  | continueStatement eos__
  | breakStatement eos__
  | blockStatement
  | functionDeclare
  ;

importStatement
  : 'import' path=StringLiteral ('as' Identifier)?
  ;

assignable
  : Identifier                        #IdentifierAssignable
  | expression '[' expression ']'     #ListIndexAssignable
  | expression '.' Identifier         #MemberAttributeAssignable
  ;

assignStatement
  : assignable (',' assignable)* '=' expression (',' expression)*
  ;

returnStatement
  : 'return' (expression (',' expression)*)?
  ;

callFunctionStatement
  : expression '(' expressionList? ')'
  ;

blockStatement
  : '{' statementList? '}'
  ;

expression
  : literal                                                   #LiteralExpression
  | Identifier                                                #IdentifierExpression
  | 'func' '(' parameterNames? ')' blockStatement             #AnonymousFunctionExpression
  | '[' expressionList? ']'                                   #ListExpression
  | '{' dictItems? '}'                                        #DictExpression
  | expression Dot Identifier                                 #MemberExpression
  | '(' expression ')'                                        #ParenthesizedExpression
  | expression '[' expression ']'                             #ListIndexExpression
  | expression '(' expressionList? ')'                        #CallFunctionExpression
  | expression op=(Multiply | Divide | Modulus) expression    #MultiplicativeExpression
  | expression op=(Plus | Minus) expression                   #AdditiveExpression
  | expression op=(LessThan | LessThanEquals | GreaterThan | GreaterThanEquals) expression 
                                                              #RelationalExpression
  | expression op=(Equals | NotEquals) expression             #EqualityExpression
  | expression BitAnd     expression                          #BitAndExpression
  | expression BitXor     expression                          #BitXorExpression
  | expression BitOr      expression                          #BitOrExpression
  | expression LogicAnd   expression                          #LogicAndExpression
  | expression LogicOr    expression                          #LogicOrExpression
  ;

expressionList
  : expression (',' expression)*
  ;

dictItems
  : Identifier ':' expression (',' Identifier ':' expression)*
  ;

literal
  : 'null' #NullLiteral
  | Lit=('true' | 'false') #BoolLiteral
  | op=(Plus | Minus)? right=FloatLiteral #FloatLiteral
  | op=(Plus | Minus)? right=IntLiteral #IntLiteral
  | StringLiteral #StringLiteral
  ;

parameterNames
  : Identifier (',' Identifier)*
  ;

elseIfList
  : ('else' 'if' cond=expression block=blockStatement)+
  ;

ifStatement
  : 'if' ifCond=expression ifBlock=blockStatement elseIfList? ('else' elseBlock=blockStatement)?
  ;

forStatement
  : 'for' expression? blockStatement
  | 'for' init=assignStatement? ';' expression? ';' inc=assignStatement? blockStatement
  ;

continueStatement: 'continue' ;

breakStatement: 'break' ;

functionDeclare
  : 'func' Identifier '(' parameterNames? ')' blockStatement
  ;

eos__
  : SemiColon
  | EOF
  | {p.lineTerminatorAhead()}?
  | {p.closeBrace()}?
  ;
