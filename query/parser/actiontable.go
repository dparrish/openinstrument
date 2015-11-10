
package parser

type(
	actionTable [numStates]actionRow
	actionRow struct {
		canRecover bool
		actions [numSymbols]action
	}
)

var actionTab = actionTable{
	actionRow{ // S0
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			shift(8),		/* percentile */
			shift(9),		/* sum */
			shift(10),		/* mean */
			shift(11),		/* median */
			shift(12),		/* min */
			shift(13),		/* max */
			shift(14),		/* stddev */
			shift(15),		/* rate */
			shift(16),		/* rate_signed */
			shift(17),		/* delta */
			shift(18),		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			shift(20),		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S1
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			accept(true),		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S2
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			reduce(1),		/* $, reduce: Query */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S3
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			reduce(2),		/* $, reduce: Query */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S4
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			reduce(3),		/* $, reduce: Query */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S5
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			shift(21),		/* ( */
			nil,		/* ) */
			shift(22),		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S6
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			shift(23),		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S7
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			reduce(18),		/* (, reduce: AggregationType */
			nil,		/* ) */
			reduce(18),		/* by, reduce: AggregationType */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S8
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			shift(24),		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S9
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			reduce(12),		/* (, reduce: AggregationType */
			nil,		/* ) */
			reduce(12),		/* by, reduce: AggregationType */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S10
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			reduce(13),		/* (, reduce: AggregationType */
			nil,		/* ) */
			reduce(13),		/* by, reduce: AggregationType */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S11
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			reduce(14),		/* (, reduce: AggregationType */
			nil,		/* ) */
			reduce(14),		/* by, reduce: AggregationType */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S12
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			reduce(15),		/* (, reduce: AggregationType */
			nil,		/* ) */
			reduce(15),		/* by, reduce: AggregationType */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S13
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			reduce(16),		/* (, reduce: AggregationType */
			nil,		/* ) */
			reduce(16),		/* by, reduce: AggregationType */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S14
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			reduce(17),		/* (, reduce: AggregationType */
			nil,		/* ) */
			reduce(17),		/* by, reduce: AggregationType */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S15
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			reduce(19),		/* (, reduce: MutationType */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S16
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			reduce(20),		/* (, reduce: MutationType */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S17
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			reduce(21),		/* (, reduce: MutationType */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S18
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			reduce(22),		/* (, reduce: MutationType */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S19
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			reduce(23),		/* $, reduce: Variable */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			shift(25),		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S20
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			reduce(26),		/* $, reduce: VarName */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			reduce(26),		/* [, reduce: VarName */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			shift(26),		/* {} */
			shift(27),		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S21
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			shift(15),		/* rate */
			shift(16),		/* rate_signed */
			shift(17),		/* delta */
			shift(18),		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			shift(33),		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S22
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			shift(34),		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S23
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			shift(36),		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S24
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			shift(38),		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S25
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			shift(39),		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S26
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			reduce(27),		/* $, reduce: VarName */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			reduce(27),		/* [, reduce: VarName */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S27
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			shift(40),		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S28
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			reduce(32),		/* ), reduce: VariableList */
			nil,		/* by */
			reduce(32),		/* ,, reduce: VariableList */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S29
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			shift(43),		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S30
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			shift(44),		/* ) */
			nil,		/* by */
			shift(45),		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S31
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			shift(46),		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S32
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			reduce(23),		/* ), reduce: Variable */
			nil,		/* by */
			reduce(23),		/* ,, reduce: Variable */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			shift(47),		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S33
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			reduce(26),		/* ), reduce: VarName */
			nil,		/* by */
			reduce(26),		/* ,, reduce: VarName */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			reduce(26),		/* [, reduce: VarName */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			shift(48),		/* {} */
			shift(49),		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S34
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			shift(51),		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S35
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			shift(52),		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S36
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			reduce(34),		/* ,, reduce: Duration */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S37
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			shift(53),		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S38
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			reduce(35),		/* ), reduce: Int64 */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S39
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			shift(54),		/* : */
			shift(55),		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S40
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			shift(56),		/* = */
			
		},

	},
	actionRow{ // S41
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			shift(57),		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			shift(58),		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S42
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			reduce(29),		/* ,, reduce: KeyValues */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			reduce(29),		/* }, reduce: KeyValues */
			nil,		/* = */
			
		},

	},
	actionRow{ // S43
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			reduce(5),		/* $, reduce: Aggregation */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S44
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			reduce(4),		/* $, reduce: Aggregation */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S45
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			shift(33),		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S46
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			shift(36),		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S47
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			shift(61),		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S48
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			reduce(27),		/* ), reduce: VarName */
			nil,		/* by */
			reduce(27),		/* ,, reduce: VarName */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			reduce(27),		/* [, reduce: VarName */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S49
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			shift(40),		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S50
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			shift(63),		/* ) */
			nil,		/* by */
			shift(64),		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S51
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			reduce(9),		/* ), reduce: LabelList */
			nil,		/* by */
			reduce(9),		/* ,, reduce: LabelList */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S52
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			shift(33),		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S53
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			reduce(11),		/* (, reduce: Percentile */
			nil,		/* ) */
			reduce(11),		/* by, reduce: Percentile */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S54
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			shift(66),		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S55
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			reduce(25),		/* $, reduce: Variable */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S56
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			shift(67),		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S57
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			shift(40),		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S58
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			reduce(28),		/* $, reduce: VarName */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			reduce(28),		/* [, reduce: VarName */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S59
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			reduce(33),		/* ), reduce: VariableList */
			nil,		/* by */
			reduce(33),		/* ,, reduce: VariableList */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S60
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			shift(69),		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S61
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			shift(70),		/* : */
			shift(71),		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S62
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			shift(57),		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			shift(72),		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S63
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			shift(73),		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S64
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			shift(74),		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S65
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			shift(75),		/* ) */
			nil,		/* by */
			shift(45),		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S66
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			shift(76),		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S67
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			reduce(31),		/* ,, reduce: KeyValue */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			reduce(31),		/* }, reduce: KeyValue */
			nil,		/* = */
			
		},

	},
	actionRow{ // S68
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			reduce(30),		/* ,, reduce: KeyValues */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			reduce(30),		/* }, reduce: KeyValues */
			nil,		/* = */
			
		},

	},
	actionRow{ // S69
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			shift(33),		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S70
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			shift(78),		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S71
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			reduce(25),		/* ), reduce: Variable */
			nil,		/* by */
			reduce(25),		/* ,, reduce: Variable */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S72
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			reduce(28),		/* ), reduce: VarName */
			nil,		/* by */
			reduce(28),		/* ,, reduce: VarName */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			reduce(28),		/* [, reduce: VarName */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S73
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			shift(15),		/* rate */
			shift(16),		/* rate_signed */
			shift(17),		/* delta */
			shift(18),		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			shift(33),		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S74
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			reduce(10),		/* ), reduce: LabelList */
			nil,		/* by */
			reduce(10),		/* ,, reduce: LabelList */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S75
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			reduce(8),		/* $, reduce: Mutation */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S76
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			reduce(24),		/* $, reduce: Variable */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S77
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			shift(81),		/* ) */
			nil,		/* by */
			shift(45),		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S78
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			shift(82),		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S79
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			shift(83),		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S80
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			shift(84),		/* ) */
			nil,		/* by */
			shift(45),		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S81
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			reduce(8),		/* ), reduce: Mutation */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S82
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			reduce(24),		/* ), reduce: Variable */
			nil,		/* by */
			reduce(24),		/* ,, reduce: Variable */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S83
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			reduce(7),		/* $, reduce: Aggregation */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S84
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			reduce(6),		/* $, reduce: Aggregation */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* percentile */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* delta */
			nil,		/* latest */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	
}

