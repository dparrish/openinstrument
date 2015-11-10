
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
		7, // Percentile
		5, // AggregationType
		6, // MutationType
		2, // Variable
		19, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S1
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		29, // Mutation
		-1, // LabelList
		-1, // Percentile
		-1, // AggregationType
		31, // MutationType
		28, // Variable
		32, // VarName
		-1, // KeyValues
		-1, // KeyValue
		30, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S22
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		35, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S24
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		37, // Int64
		

	},
	gotoRow{ // S25
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S27
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
		-1, // Variable
		-1, // VarName
		41, // KeyValues
		42, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S28
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		50, // LabelList
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // LabelList
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
		59, // Variable
		32, // VarName
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		60, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S47
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
		-1, // Variable
		-1, // VarName
		62, // KeyValues
		42, // KeyValue
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
		28, // Variable
		32, // VarName
		-1, // KeyValues
		-1, // KeyValue
		65, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S53
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S54
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S56
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		68, // KeyValue
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S63
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
		28, // Variable
		32, // VarName
		-1, // KeyValues
		-1, // KeyValue
		77, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S70
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // LabelList
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		79, // Mutation
		-1, // LabelList
		-1, // Percentile
		-1, // AggregationType
		31, // MutationType
		28, // Variable
		32, // VarName
		-1, // KeyValues
		-1, // KeyValue
		80, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	gotoRow{ // S74
		
		-1, // S'
		-1, // Query
		-1, // Aggregation
		-1, // Mutation
		-1, // LabelList
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
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
		-1, // Percentile
		-1, // AggregationType
		-1, // MutationType
		-1, // Variable
		-1, // VarName
		-1, // KeyValues
		-1, // KeyValue
		-1, // VariableList
		-1, // Duration
		-1, // Int64
		

	},
	
}
