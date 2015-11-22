
/*
*/
package parser

const numNTSymbols = 15
type(
	gotoTable [numStates]gotoRow
	gotoRow	[numNTSymbols] int
)

var gotoTab = gotoTable{
	gotoRow{ // S0
		
		-1, // S'
		1, // Query
		4, // Aggregation
		3, // Mutation
		-1, // LabelList
		5, // AggregationType
		6, // MutationType
		13, // Percentile
		20, // Variable
		21, // VarName
		-1, // KeyValues
		-1, // KeyValue
		2, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S1
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S2
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S3
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S4
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S5
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S6
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S7
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S8
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S9
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S10
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S11
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S12
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S13
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S14
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S15
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S16
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S17
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S18
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S19
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S20
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S21
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S22
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S23
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		31, // Variable
		21, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S24
		
		-1, // S'
		32, // Query
		35, // Aggregation
		34, // Mutation
		-1, // LabelList
		36, // AggregationType
		37, // MutationType
		13, // Percentile
		38, // Variable
		39, // VarName
		-1, // KeyValues
		-1, // KeyValue
		33, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S25
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S26
		
		-1, // S'
		42, // Query
		35, // Aggregation
		34, // Mutation
		-1, // LabelList
		36, // AggregationType
		37, // MutationType
		13, // Percentile
		38, // Variable
		39, // VarName
		-1, // KeyValues
		-1, // KeyValue
		33, // VariableList
		43, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S27
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		45, // Int64
		

	},
	gotoRow{ // S28
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S29
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S30
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		49, // KeyValues
		50, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S31
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S32
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S33
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S34
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S35
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S36
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S37
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S38
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S39
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S40
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S41
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		59, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S42
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S43
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S44
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S45
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S46
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S47
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S48
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S49
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S50
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S51
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S52
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		69, // Variable
		39, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S53
		
		-1, // S'
		70, // Query
		35, // Aggregation
		34, // Mutation
		-1, // LabelList
		36, // AggregationType
		37, // MutationType
		13, // Percentile
		38, // Variable
		39, // VarName
		-1, // KeyValues
		-1, // KeyValue
		33, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S54
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S55
		
		-1, // S'
		72, // Query
		35, // Aggregation
		34, // Mutation
		-1, // LabelList
		36, // AggregationType
		37, // MutationType
		13, // Percentile
		38, // Variable
		39, // VarName
		-1, // KeyValues
		-1, // KeyValue
		33, // VariableList
		73, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S56
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S57
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S58
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		75, // KeyValues
		50, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S59
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S60
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S61
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S62
		
		-1, // S'
		78, // Query
		35, // Aggregation
		34, // Mutation
		-1, // LabelList
		36, // AggregationType
		37, // MutationType
		13, // Percentile
		38, // Variable
		39, // VarName
		-1, // KeyValues
		-1, // KeyValue
		33, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S63
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S64
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S65
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S66
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S67
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		81, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S68
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S69
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S70
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S71
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		83, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S72
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S73
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S74
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S75
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S76
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S77
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S78
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S79
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S80
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S81
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S82
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S83
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S84
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S85
		
		-1, // S'
		94, // Query
		35, // Aggregation
		34, // Mutation
		-1, // LabelList
		36, // AggregationType
		37, // MutationType
		13, // Percentile
		38, // Variable
		39, // VarName
		-1, // KeyValues
		-1, // KeyValue
		33, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S86
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S87
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S88
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S89
		
		-1, // S'
		96, // Query
		35, // Aggregation
		34, // Mutation
		-1, // LabelList
		36, // AggregationType
		37, // MutationType
		13, // Percentile
		38, // Variable
		39, // VarName
		-1, // KeyValues
		-1, // KeyValue
		33, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S90
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S91
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S92
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S93
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S94
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S95
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S96
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S97
		
		-1, // S'
		101, // Query
		35, // Aggregation
		34, // Mutation
		-1, // LabelList
		36, // AggregationType
		37, // MutationType
		13, // Percentile
		38, // Variable
		39, // VarName
		-1, // KeyValues
		-1, // KeyValue
		33, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S98
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S99
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S100
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S101
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S102
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // AggregationType
		-1, // MutationType
		-1, // Percentile
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	
}
