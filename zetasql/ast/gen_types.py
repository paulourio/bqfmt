from dataclasses import dataclass, field
from enum import Enum
from typing import Any, Dict, List, Optional
import os
import textwrap

import jinja2


class FieldLoader(Enum):

    REQUIRED = 1
    OPTIONAL = 2
    REPEATED = 3


@dataclass
class Field:

    name: str
    type: str
    field_loader: FieldLoader = FieldLoader.OPTIONAL
    comment: Optional[str] = None
    init: bool = True  # Whether the field should be required in New()

    @property
    def arg_name(self):
        arg_name = self.name.lower()
        if arg_name == 'type':
            arg_name = 'typ'
        return arg_name

    @property
    def var_name(self):
        if len(self.name) <= 3:
            return 'var' + self.name
        return self.name.lower()[:3]

    @property
    def typed_var_name(self):
        return 'typed' + self.name


class TreeGenerator:

    def __init__(self):
        self.nodes: List[Dict[str, Any]] = []
        self.node_map: Dict[str, Dict[str, Any]] = {}

    def add_node(self, name: str, composition: str,
                 fields: Optional[List[Field]] = None,
                 comment: Optional[str] = None,
                 custom_methods: List[str] = None,
                 attribs: List[str] = None) -> None:
        if fields is None:
            fields = []
        if custom_methods is None:
            custom_methods = []
        has_repeated = False
        for f in fields:
            if f.field_loader == FieldLoader.REPEATED:
                if has_repeated:
                    raise RuntimeError(
                        f'only one repeated allowed in {f.name}')
                has_repeated = True
        node = {
            'name': name,
            'composition': composition,
            'fields': fields,
            'comment': comment,
            'custom_methods': custom_methods,
            'attribs': attribs,
        }
        assert name not in self.node_map, f'duplicate node name {name}'
        self.nodes.append(node)

    def generate(self, output_path: str) -> None:
        """Materialize the template to generate the output files."""
        env = jinja2.Environment(
            undefined=jinja2.StrictUndefined,
            loader=jinja2.FileSystemLoader(''),
        )
        env.globals['make_comment'] = make_comment
        context = {
            'nodes': self.nodes,
            'FieldLoader': FieldLoader,
        }
        template = env.get_template('types_generated.go.j2')
        oname = os.path.join(output_path, 'types_generated.go')
        print('Writing', oname)
        with open(oname, 'w') as out:
            out.write(template.render(context))


def make_comment(text: str, indent: int, width: int) -> str:
    assert isinstance(text, str)
    prefix = '\n' + ' '*indent + '// '
    return prefix.join(textwrap.wrap('// ' + text, width - 3))


def main():
    gen = TreeGenerator()

    gen.add_node(
        name='QueryStatement',
        composition='Statement',
        comment='QueryStatement represents a single query statement.',
        fields=[
            Field('Query', '*Query', FieldLoader.REQUIRED)
        ])

    gen.add_node(
        name='Query',
        composition='QueryExpression',
        fields=[
            Field(
                'WithClause',
                '*WithClause',
                comment='WithClause is the WITH clause wrapping this query.'),
            Field(
                'QueryExpr',
                'QueryExpressionHandler',
                field_loader=FieldLoader.REQUIRED,
                comment=(
                    'QueryExpr can be a single Select, or a more complex '
                    'structure composed out of nodes like SetOperation '
                    'and Query.'
                )),
            Field(
                'OrderBy',
                '*OrderBy',
                comment=(
                    'OrderBy applies, if present, to the result of QueryExpr.'
                )),
            Field(
                'LimitOffset',
                '*LimitOffset',
                comment=(
                    'LimitOffset applies, if present, after the result '
                    'of QueryExpr and OrderBy.'
                ))
        ])

    gen.add_node(
        name='Select',
        composition='QueryExpression',
        custom_methods=['DebugString'],
        fields=[
            Field('Distinct', 'bool'),
            Field('SelectAs', '*SelectAs'),
            Field('SelectList', '*SelectList', FieldLoader.REQUIRED),
            Field('FromClause', '*FromClause'),
            Field('WhereClause', '*WhereClause'),
            Field('GroupBy', '*GroupBy'),
            Field('Having', '*Having'),
            Field('Qualify', '*Qualify'),
            Field('WindowClause', '*WindowClause'),
        ])

    gen.add_node(
        name='SelectList',
        composition='Node',
        fields=[
            Field('Columns', '*SelectColumn', FieldLoader.REPEATED)
        ])

    gen.add_node(
        name='SelectColumn',
        composition='Node',
        fields=[
            Field('Expression', 'ExpressionHandler', FieldLoader.REQUIRED),
            Field('Alias', '*Alias'),
        ])

    gen.add_node(name='IntLiteral', composition='Leaf'),

    gen.add_node(
        name='Identifier',
        composition='Expression',
        custom_methods=['DebugString'],
        fields=[
            Field('IDString', 'string')
        ])

    gen.add_node(
        name='Alias',
        composition='Node',
        fields=[
            Field('Identifier', '*Identifier', FieldLoader.REQUIRED)
        ])

    gen.add_node(
        name='PathExpression',
        composition='Expression',
        comment=(
            'PathExpression is used for dotted identifier paths only, '
            'not dotting into arbitrary expressions (see DotIdentifier).'
        ),
        fields=[
            Field('Names', '*Identifier', FieldLoader.REPEATED)
        ])

    gen.add_node(
        name='TablePathExpression',
        composition='TableExpression',
        comment=(
            'TablePathExpression is a table expression than introduce '
            'a single scan, referenced by a path expression or UNNEST, '
            'can optionally have aliases.  Exactly one of PathExpr '
            'or UnnestExpr must be non nil.'
        ),
        fields=[
            Field('PathExpr', '*PathExpression'),
            Field('UnnestExpr', '*UnnestExpression'),
            Field('Alias', '*Alias'),
        ])

    gen.add_node(
        name='FromClause',
        composition='Node',
        fields=[
            Field(
                'TableExpression',
                'TableExpressionHandler',
                field_loader=FieldLoader.REQUIRED,
                comment=(
                    'TableExpression has exactly one table expression '
                    'child.  If the FROM clause has commas, they will '
                    'expressed as a tree of Join nodes with JoinType=Comma.'
                )),
        ])

    gen.add_node(
        name='WhereClause',
        composition='Node',
        fields=[
            Field('Expression', 'ExpressionHandler', FieldLoader.REQUIRED)
        ])

    gen.add_node(
        name='BooleanLiteral',
        composition='Leaf',
        fields=[
            Field('Value', 'bool')
        ])

    gen.add_node(
        name='AndExpr',
        composition='Expression',
        fields=[
            Field('Conjuncts', 'ExpressionHandler', FieldLoader.REPEATED),
        ])

    gen.add_node(
        name='BinaryExpression',
        composition='Expression',
        custom_methods=['DebugString'],
        fields=[
            Field('Op', 'BinaryOp'),
            Field('LHS', 'ExpressionHandler', FieldLoader.REQUIRED),
            Field('RHS', 'ExpressionHandler', FieldLoader.REQUIRED),
            Field(
                'IsNot',
                'bool',
                comment=(
                    'IsNot indicates whether the binary operator has a '
                    'preceding NOT to it.  For NOT LIKE and IS NOT.')),
        ])

    gen.add_node(
        name='StringLiteral',
        composition='Leaf',
        fields=[Field('StringValue', 'string', init=False)])

    gen.add_node(
        name='Star',
        composition='Leaf')

    gen.add_node(
        name='OrExpr',
        composition='Expression',
        fields=[
            Field('Disjuncts', 'ExpressionHandler', FieldLoader.REPEATED)
        ])

    gen.add_node(
        name='GroupingItem',
        composition='Node',
        comment=(
            'GroupingItem represents a grouping item, which is either '
            'an expression (a regular group by key) or a rollup list. '
            'Exactly one of Expression and Rollup will be non-nil.'),
        fields=[
            Field('Expression', 'ExpressionHandler'),
            Field('Rollup', '*Rollup')
        ])

    gen.add_node(
        name='GroupBy',
        composition='Node',
        fields=[
            Field('GroupingItems', '*GroupingItem', FieldLoader.REPEATED)
        ])

    gen.add_node(
        name='OrderingExpression',
        composition='Node',
        fields=[
            Field('Expression', 'ExpressionHandler', FieldLoader.REQUIRED),
            Field('NullOrder', '*NullOrder'),
            Field('OrderingSpec', 'OrderingSpec'),
        ])

    gen.add_node(
        name='OrderBy',
        composition='Node',
        fields=[
            Field(
                'OrderingExpression',
                '*OrderingExpression',
                FieldLoader.REPEATED)
        ])

    gen.add_node(
        name='LimitOffset',
        composition='Node',
        fields=[
            Field(
                'Limit',
                'ExpressionHandler',
                FieldLoader.REQUIRED,
                comment='Limit is the LIMIT value, never nil.'),
            Field(
                'Offset',
                'ExpressionHandler',
                comment='Offset is the optional OFFSET value, or nil.'),
        ])

    gen.add_node(
        name='FloatLiteral',
        composition='Leaf')

    gen.add_node(
        name='NullLiteral',
        composition='Leaf')

    gen.add_node(
        name='OnClause',
        composition='Node',
        fields=[
            Field('Expression', 'ExpressionHandler', FieldLoader.REQUIRED),
        ])

    gen.add_node(
        name='WithClauseEntry',
        composition='Node',
        fields=[
            Field('Alias', '*Identifier', FieldLoader.REQUIRED),
            Field('Query', '*Query', FieldLoader.REQUIRED),
        ])

    gen.add_node(
        name='Join',
        composition='TableExpression',
        comment=(
            'Join can introduce multiple scans and cannot have aliases. '
            'It can also represent a JOIN with a list of consecutive '
            'ON/USING clauses.'
        ),
        fields=[
            Field('LHS', 'TableExpressionHandler', FieldLoader.REQUIRED),
            Field('RHS', 'TableExpressionHandler', FieldLoader.REQUIRED),
            Field('ClauseList', '*OnOrUsingClauseList'),
            Field('JoinType', 'JoinType'),
        ])

    gen.add_node(
        name='WithClause',
        composition='Node',
        fields=[
            Field('With', '*WithClauseEntry', FieldLoader.REPEATED),
        ])

    gen.add_node(
        name='Having',
        composition='Node')

    gen.add_node(
        name='NamedType',
        composition='Type',
        fields=[
            Field('TypeName', '*PathExpression', FieldLoader.REQUIRED),
            Field('TypeParameters', '*TypeParameterList')
        ])

    gen.add_node(
        name='ArrayType',
        composition='Type',
        fields=[
            Field('ElementType', 'TypeHandler', FieldLoader.REQUIRED),
            Field('TypeParameters', '*TypeParameterList'),
        ])

    gen.add_node(
        name='StructField',
        composition='Node',
        fields=[
            Field(
                'Name',
                '*Identifier',
                comment=(
                    'Name will be nil for anonymous fields '
                    'like in STRUCT<int, string>.')),
        ])

    gen.add_node(
        name='StructType',
        composition='Type',
        fields=[
            Field('StructFields', '*StructField', FieldLoader.REQUIRED),
            Field('TypeParameterList', '*TypeParameterList'),
        ])

    gen.add_node(
        name='CastExpression',
        composition='Expression',
        fields=[
            Field('Expr', 'ExpressionHandler', FieldLoader.REQUIRED),
            Field('Type', 'TypeHandler', FieldLoader.REQUIRED),
            Field('Format', '*FormatClause'),
            Field('IsSafeCast', 'bool'),
        ])

    gen.add_node(
        name='SelectAs',
        composition='Node',
        comment=(
            'SelectAs represents a SELECT with AS clause giving it an '
            'output type. Exactly one of SELECT AS STRUCT, SELECT AS VALUE, '
            'SELECT AS <TypeName> is present.'),
        fields=[
            Field('TypeName', '*PathExpression'),
            Field('AsMode', 'AsMode'),
        ])

    gen.add_node(
        name='Rollup',
        composition='Node',
        fields=[
            Field('Expressions', 'ExpressionHandler', FieldLoader.REPEATED),
        ])

    gen.add_node(
        name='FunctionCall',
        composition='Expression',
        fields=[
            Field('Function', '*PathExpression', FieldLoader.REQUIRED),
            Field('Arguments', 'ExpressionHandler', FieldLoader.REPEATED),
            Field(
                'OrderBy',
                '*OrderBy',
                comment=(
                    'OrderBy is set when the function is called with '
                    'FUNC(args ORDER BY cols).')),
            Field(
                'LimitOffset',
                '*LimitOffset',
                comment=(
                    'LimitOffset is set when the function is called '
                    'with FUNC(args LIMIT n).')),
            Field(
                'NullHandlingModifier',
                'NullHandlingModifier',
                comment=(
                    'NullHandlingModifier is set when the function is '
                    'called with FUNC(args {IGNORE|RESPECT} NULLS).')),
            Field(
                'Distinct',
                'bool',
                comment=(
                    'Distinct is true when the function is called with '
                    'FUNC(DISTINCT args).')),
        ])

    gen.add_node(
        name='ArrayConstructor',
        composition='Expression',
        fields=[
            Field(
                'Type',
                '*ArrayType',
                comment=(
                    'Type may be nil, depending on whether the array is '
                    'constructed through ARRAY<type>[...] syntax or '
                    'ARRAY[...] or [...].')),
            Field(
                'Elements',
                'ExpressionHandler',
                FieldLoader.REPEATED),
        ])

    gen.add_node(
        name='StructConstructorArg',
        composition='Node',
        fields=[
            Field('Expression', 'ExpressionHandler', FieldLoader.REQUIRED),
            Field('Alias', '*Alias', FieldLoader.REQUIRED),
        ])

    gen.add_node(
        name='StructConstructorWithParens',
        composition='Node',
        comment=(
            'StructConstructorWithParens is resulted from structs '
            'constructed with (expr, expr, ...) with at least two '
            'expressions.'),
        fields=[
            Field(
                'FieldExpressions',
                'ExpressionHandler',
                FieldLoader.REPEATED)
        ])

    gen.add_node(
        name='StructConstructorWithKeyword',
        composition='Expression',
        comment=(
            'StructConstructorWithKeyword is resulted from structs '
            'constructed with STRUCT(expr [AS alias], ...) or '
            'STRUCT<...>(expr [AS alias], ...). Both forms support '
            'empty field lists.  The StructType is non-nil when '
            'the type is explicitly defined.'),
        fields=[
            Field('StructType', '*StructType'),
            Field('Fields', '*StructConstructorArg', FieldLoader.REPEATED),
        ])

    gen.add_node(
        name='InExpression',
        composition='Expression',
        comment=(
            'InExpression is resulted from expr IN (expr, expr, ...), '
            'expr IN UNNEST(...), and expr IN (query). Exactly one of '
            'InList, Query, or UnnestExpr is present.'),
        fields=[
            Field('LHS', 'ExpressionHandler', FieldLoader.REQUIRED),
            Field('InList', '*InList'),
            Field('Query', '*Query'),
            Field('UnnestExpr', '*UnnestExpression'),
            Field(
                'IsNot',
                'bool',
                comment=(
                    'IsNot signifies whether the IN operator as a '
                    'preceding NOT to it.')),
        ])

    gen.add_node(
        name='InList',
        composition='Node',
        comment=(
            'InList is shared with the IN operator and LIKE ANY/SOME/ALL.'),
        fields=[
            Field(
                'List',
                'ExpressionHandler',
                FieldLoader.REPEATED,
                comment=(
                    'List contains the expressions present in the '
                    'InList node.'
                ))
        ])

    gen.add_node(
        name='BetweenExpression',
        composition='Expression',
        comment=(
            'BetweenExpression is resulted through <LHS> BETWEEN <Low> '
            'AND <High>.'),
        fields=[
            Field('LHS', 'ExpressionHandler', FieldLoader.REQUIRED),
            Field('Low', 'ExpressionHandler', FieldLoader.REQUIRED),
            Field('High', 'ExpressionHandler', FieldLoader.REQUIRED),
            Field(
                'IsNot',
                'bool',
                comment=(
                    'IsNot signifies whether the BETWEEN operator '
                    'has a preceding NOT to it.'
                )),
        ])

    gen.add_node(name='NumericLiteral', composition='Leaf')

    gen.add_node(name='BigNumericLiteral', composition='Leaf')

    gen.add_node(name='BytesLiteral', composition='Leaf')

    gen.add_node(
        name='DateOrTimeLiteral',
        composition='Expression',
        fields=[
            Field('StringLiteral', '*StringLiteral', FieldLoader.REQUIRED),
            Field('TypeKind', 'TypeKind')
        ])

    gen.add_node(
        name='CaseValueExpression',
        composition='Expression',
        fields=[
            Field('Arguments', 'ExpressionHandler', FieldLoader.REPEATED),
        ])

    gen.add_node(
        name='CaseNoValueExpression',
        composition='Expression',
        fields=[
            Field('Arguments', 'ExpressionHandler', FieldLoader.REPEATED),
        ]
    )

    gen.add_node(
        name='ArrayElement',
        composition='Expression',  # GeneralizedPathExpression
        fields=[
            Field('Array', 'ExpressionHandler', FieldLoader.REQUIRED),
            Field('Position', 'ExpressionHandler', FieldLoader.REQUIRED)
        ])

    gen.add_node(
        name='BitwiseShiftExpression',
        composition='Expression',
        fields=[
            Field('LHS', 'ExpressionHandler', FieldLoader.REQUIRED),
            Field('RHS', 'ExpressionHandler', FieldLoader.REQUIRED),
            Field(
                'IsLeftShift',
                'bool',
                comment=(
                    'IsLeftShift signifies whether the bitwise shift '
                    'is of left shift type "<<" or right shift type ">>".')),
        ])

    gen.add_node(
        name='DotGeneralizedField',
        composition='Expression',  # GeneralizedPathExpression
        comment=(
            'DotGeneralizedField is a generalized form of extracting '
            'a field from an expression. It uses a parenthesized '
            'PathExpression instead of a single identifier ot select a '
            'field.'),
        fields=[
            Field('Expr', 'ExpressionHandler', FieldLoader.REQUIRED),
            Field('Path', '*PathExpression', FieldLoader.REQUIRED),
        ])

    gen.add_node(
        name='DotIdentifier',
        composition='Expression',  # GeneralizedPathExpression
        comment=(
            'DotIdentifier is used for using dot to extract a field '
            'from an arbitrary expression. Is cases where we know the '
            'left side is always an identifier path, we use '
            'PathExpression instead.'),
        fields=[
            Field('Expr', 'ExpressionHandler', FieldLoader.REQUIRED),
            Field('Path', '*PathExpression', FieldLoader.REQUIRED),
        ])

    gen.add_node(
        name='DotStar',
        composition='Expression',
        fields=[
            Field('Expr', 'ExpressionHandler', FieldLoader.REQUIRED),
        ])

    gen.add_node(
        name='DotStarWithModifiers',
        composition='Expression',
        comment=(
            'DotStarWithModifiers is an expression constructed through '
            'SELECT x.* EXCEPT (...) REPLACE (...).'),
        fields=[
            Field('Expr', 'ExpressionHandler', FieldLoader.REQUIRED),
            Field('Modifiers', '*StarModifiers', FieldLoader.REQUIRED),
        ])

    gen.add_node(
        name='ExpressionSubquery',
        composition='Expression',
        comment=(
            'ExpressionSubquery is a subquery in an expression. '
            '(Not in the FROM clause.)'),
        fields=[
            Field('Query', '*Query', FieldLoader.REQUIRED),
            Field(
                'Modifier',
                'SubqueryModifier',
                FieldLoader.REQUIRED,
                comment=(
                    'Modifier is the syntactic modifier on this '
                    'expression subquery.')),
        ])

    gen.add_node(
        name='ExtractExpression',
        composition='Expression',
        comment=(
            'ExtractExpression is resulted from '
            'EXTRACT(<LHS> FROM <RHS> <TimeZone>).'
        ),
        fields=[
            Field('LHS', 'ExpressionHandler', FieldLoader.REQUIRED),
            Field('RHS', 'ExpressionHandler', FieldLoader.REQUIRED),
            Field('TimeZone', 'ExpressionHandler'),
        ])

    gen.add_node(
        name='IntervalExpression',
        composition='Expression',
        fields=[
            Field('IntervalValue', 'ExpressionHandler', FieldLoader.REQUIRED),
            Field('DatePartName', 'ExpressionHandler', FieldLoader.REQUIRED),
            Field('DatePartNameTo', 'ExpressionHandler'),
        ])

    gen.add_node(
        name='NullOrder',
        composition='Node',
        fields=[
            Field('NullsFirst', 'bool')
        ])

    gen.add_node(
        name='OnOrUsingClauseList',
        composition='Node',
        fields=[
            Field(
                'List',
                'NodeHandler',
                FieldLoader.REPEATED,
                comment=(
                    'List is a list of OnClause and UsingClause elements.'
                ))
        ])

    gen.add_node(
        name='ParenthesizedJoin',
        composition='TableExpression',
        fields=[
            Field('Join', '*Join', FieldLoader.REQUIRED),
            Field('SampleClause', 'SampleClause'),
        ])

    gen.add_node(
        name='PartitionBy',
        composition='Node',
        fields=[
            Field(
                'PartitioningExpressions',
                'ExpressionHandler',
                FieldLoader.REPEATED)
        ])

    gen.add_node(
        name='SetOperation',
        composition='QueryExpression',
        fields=[
            Field('Inputs', 'QueryExpressionHandler', FieldLoader.REPEATED),
            Field('OpType', 'SetOp'),
            Field('Distinct', 'bool'),
        ])

    gen.add_node(
        name='StarExceptList',
        composition='Node',
        fields=[
            Field('Identifiers', '*Identifier', FieldLoader.REPEATED),
        ])

    gen.add_node(
        name='StarModifiers',
        composition='Node',
        comment=(
            'StarModifiers is resulted from SELECT * EXCEPT (...) '
            'REPLACE (...).'
        ),
        fields=[
            Field('ExceptList', '*StarExceptList'),
            Field('ReplaceItems', '*StarReplaceItem', FieldLoader.REPEATED),
        ])

    gen.add_node(
        name='StarReplaceItem',
        composition='Node',
        fields=[
            Field('Expression', 'ExpressionHandler', FieldLoader.REQUIRED),
            Field('Alias', '*Identifier', FieldLoader.REQUIRED),
        ])

    gen.add_node(
        name='StarWithModifiers',
        composition='Expression',
        comment=(
            'StarModifiers is resulted from SELECT * EXCEPT (...) '
            'REPLACE (...).'),
        fields=[
            Field('Modifiers', '*StarModifiers', FieldLoader.REQUIRED),
        ])

    gen.add_node(
        name='TableSubquery',
        composition='TableExpression',
        comment=(
            'TableSubquery contains the table subquery, which can contain '
            'either a PivotClause or an UnpivotClause.'),
        fields=[
            Field('Subquery', '*Query', FieldLoader.REQUIRED),
            Field('Alias', '*Alias'),
            # Field('PivotClause', '*PivotClause'),
            # Field('UnpivotClause', '*UnpivotClause'),
            Field('SampleClause', '*SampleClause'),
        ])

    gen.add_node(
        name='UnaryExpression',
        composition='Expression',
        fields=[
            Field('Operand', 'ExpressionHandler', FieldLoader.REQUIRED),
            Field('Op', 'UnaryOp'),
        ])

    gen.add_node(
        name='UnnestExpression',
        composition='Node',
        fields=[
            Field('Expression', 'ExpressionHandler', FieldLoader.REQUIRED),
        ])

    gen.add_node(
        name='WindowClause',
        composition='Node',
        fields=[
            Field('Windows', '*WindowDefinition', FieldLoader.REPEATED),
        ])

    gen.add_node(
        name='WindowDefinition',
        composition='Node',
        fields=[
            Field('Name', '*Identifier', FieldLoader.REQUIRED),
            Field('WindowSpec', '*WindowSpecification', FieldLoader.REQUIRED),
        ])

    gen.add_node(
        name='WindowFrame',
        composition='Node',
        fields=[
            Field('StartExpr', '*WindowFrameExpression', FieldLoader.REQUIRED),
            Field('EndExpr', '*WindowFrameExpression', FieldLoader.REQUIRED),
            Field('FrameUnit', 'FrameUnit', FieldLoader.REQUIRED),
        ])

    gen.add_node(
        name='WindowFrameExpression',
        composition='Node',
        fields=[
            Field(
                'Expression',
                'ExpressionHandler',
                FieldLoader.REQUIRED,
                comment=(
                    'Expression specifies the boundary as a logical or '
                    'physical offset to current row. It is present when '
                    'BoundaryType is OffsetPreceding or OffsetFollowing.')),
            Field('BoundaryType', 'BoundaryType'),
        ])

    gen.add_node(
        name='LikeExpression',
        composition='Expression',
        fields=[
            Field('LHS', 'ExpressionHandler', FieldLoader.REQUIRED),
            Field('InList', '*InList', FieldLoader.REQUIRED),
            Field('IsNot', 'bool'),
        ])

    gen.add_node(
        name='WindowSpecification',
        composition='Node',
        fields=[
            Field('BaseWindowName', '*Identifier'),
            Field('PartitionBy', '*PartitionBy'),
            Field('OrderBy', '*OrderBy'),
            Field('WindowFrame', '*WindowFrame'),
        ])

    gen.add_node(
        name='WithOffset',
        composition='Node',
        fields=[
            Field('Alias', '*Alias'),
        ])

    gen.add_node(
        name='TypeParameterList',
        composition='Node',
        fields=[
            Field('Parameters', 'LeafHandler', FieldLoader.REPEATED),
        ])

    gen.add_node(
        name='SampleClause',
        composition='Node',
        fields=[
            Field('SampleMethod', '*Identifier', FieldLoader.REQUIRED),
            Field('SampleSize', '*SampleSize', FieldLoader.REQUIRED),
            Field('SampleSuffix', '*SampleSuffix'),
        ])

    gen.add_node(
        name='SampleSize',
        composition='Node',
        fields=[
            Field('Size', 'ExpressionHandler', FieldLoader.REQUIRED),
            Field('PartitionBy', '*PartitionBy'),
            Field('Unit', 'SampleSizeUnit'),
        ])

    gen.add_node(
        name='SampleSuffix',
        composition='Node',
        fields=[
            Field('Weight', '*WithWeight'),
            Field('Repeat', '*RepeatableClause'),
        ])

    gen.add_node(
        name='WithWeight',
        composition='Node',
        fields=[
            Field('Alias', '*Alias'),
        ])

    gen.add_node(
        name='RepeatableClause',
        composition='Node',
        fields=[
            Field('Argument', 'ExpressionHandler', FieldLoader.REQUIRED),
        ])

    gen.add_node(
        name='Qualify',
        composition='Node',
        fields=[
            Field('Expression', 'ExpressionHandler', FieldLoader.REQUIRED)
        ])

    gen.add_node(
        name='FormatClause',
        composition='Node',
        fields=[
            Field('Format', 'ExpressionHandler', FieldLoader.REQUIRED),
            Field('TimeZoneExpr', 'ExpressionHandler'),
        ])


    gen.generate('.')


if __name__ == '__main__':
    main()
