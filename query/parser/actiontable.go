
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
			shift(7),		/* sum */
			shift(8),		/* mean */
			shift(9),		/* median */
			shift(10),		/* min */
			shift(11),		/* max */
			shift(12),		/* stddev */
			shift(14),		/* rate */
			shift(15),		/* rate_signed */
			shift(16),		/* latest */
			shift(17),		/* regularize */
			shift(18),		/* interpolate */
			shift(19),		/* percentile */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			shift(22),		/* varname */
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
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			shift(23),		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			shift(24),		/* ( */
			nil,		/* ) */
			shift(25),		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			shift(26),		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			reduce(10),		/* (, reduce: AggregationType */
			nil,		/* ) */
			reduce(10),		/* by, reduce: AggregationType */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			reduce(11),		/* (, reduce: AggregationType */
			nil,		/* ) */
			reduce(11),		/* by, reduce: AggregationType */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			reduce(17),		/* (, reduce: MutationType */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			reduce(18),		/* (, reduce: MutationType */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			reduce(19),		/* (, reduce: MutationType */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			reduce(20),		/* (, reduce: MutationType */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			reduce(21),		/* (, reduce: MutationType */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			nil,		/* $ */
			shift(27),		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S20
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			reduce(32),		/* $, reduce: VariableList */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			reduce(32),		/* ,, reduce: VariableList */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S21
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			reduce(23),		/* $, reduce: Variable */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			reduce(23),		/* ,, reduce: Variable */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
			shift(28),		/* [ */
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
	actionRow{ // S22
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			reduce(26),		/* $, reduce: VarName */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			reduce(26),		/* ,, reduce: VarName */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
			reduce(26),		/* [, reduce: VarName */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			shift(29),		/* {} */
			shift(30),		/* { */
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
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			shift(22),		/* varname */
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
			shift(7),		/* sum */
			shift(8),		/* mean */
			shift(9),		/* median */
			shift(10),		/* min */
			shift(11),		/* max */
			shift(12),		/* stddev */
			shift(14),		/* rate */
			shift(15),		/* rate_signed */
			shift(16),		/* latest */
			shift(17),		/* regularize */
			shift(18),		/* interpolate */
			shift(19),		/* percentile */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			shift(40),		/* varname */
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
			shift(41),		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S26
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			shift(44),		/* varlabel */
			shift(7),		/* sum */
			shift(8),		/* mean */
			shift(9),		/* median */
			shift(10),		/* min */
			shift(11),		/* max */
			shift(12),		/* stddev */
			shift(14),		/* rate */
			shift(15),		/* rate_signed */
			shift(16),		/* latest */
			shift(17),		/* regularize */
			shift(18),		/* interpolate */
			shift(19),		/* percentile */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			shift(40),		/* varname */
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
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
			nil,		/* [ */
			shift(46),		/* int64 */
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
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
			nil,		/* [ */
			shift(47),		/* int64 */
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
			reduce(27),		/* $, reduce: VarName */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			reduce(27),		/* ,, reduce: VarName */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S30
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			shift(48),		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			reduce(33),		/* $, reduce: VariableList */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			reduce(33),		/* ,, reduce: VariableList */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			shift(51),		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S33
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			reduce(1),		/* ), reduce: Query */
			nil,		/* by */
			shift(52),		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S34
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			reduce(2),		/* ), reduce: Query */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			reduce(3),		/* ), reduce: Query */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			shift(53),		/* ( */
			nil,		/* ) */
			shift(54),		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			shift(55),		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			reduce(32),		/* ), reduce: VariableList */
			nil,		/* by */
			reduce(32),		/* ,, reduce: VariableList */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			reduce(23),		/* ), reduce: Variable */
			nil,		/* by */
			reduce(23),		/* ,, reduce: Variable */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
			shift(56),		/* [ */
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
	actionRow{ // S40
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			reduce(26),		/* ), reduce: VarName */
			nil,		/* by */
			reduce(26),		/* ,, reduce: VarName */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
			reduce(26),		/* [, reduce: VarName */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			shift(57),		/* {} */
			shift(58),		/* { */
			nil,		/* } */
			nil,		/* = */
			
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
			nil,		/* , */
			shift(60),		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S42
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			shift(61),		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S43
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			shift(62),		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			reduce(34),		/* ,, reduce: Duration */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			shift(63),		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S46
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			reduce(35),		/* ), reduce: Int64 */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
			nil,		/* [ */
			nil,		/* int64 */
			shift(64),		/* : */
			shift(65),		/* ] */
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
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			shift(66),		/* = */
			
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
			shift(67),		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			shift(68),		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S50
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			reduce(29),		/* ,, reduce: KeyValues */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S51
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			reduce(4),		/* $, reduce: Aggregation */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			shift(40),		/* varname */
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
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			shift(7),		/* sum */
			shift(8),		/* mean */
			shift(9),		/* median */
			shift(10),		/* min */
			shift(11),		/* max */
			shift(12),		/* stddev */
			shift(14),		/* rate */
			shift(15),		/* rate_signed */
			shift(16),		/* latest */
			shift(17),		/* regularize */
			shift(18),		/* interpolate */
			shift(19),		/* percentile */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			shift(40),		/* varname */
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
			shift(71),		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S55
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			shift(44),		/* varlabel */
			shift(7),		/* sum */
			shift(8),		/* mean */
			shift(9),		/* median */
			shift(10),		/* min */
			shift(11),		/* max */
			shift(12),		/* stddev */
			shift(14),		/* rate */
			shift(15),		/* rate_signed */
			shift(16),		/* latest */
			shift(17),		/* regularize */
			shift(18),		/* interpolate */
			shift(19),		/* percentile */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			shift(40),		/* varname */
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
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
			nil,		/* [ */
			shift(74),		/* int64 */
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
			reduce(27),		/* ), reduce: VarName */
			nil,		/* by */
			reduce(27),		/* ,, reduce: VarName */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S58
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			shift(48),		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S59
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			shift(76),		/* ) */
			nil,		/* by */
			shift(77),		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			reduce(8),		/* ), reduce: LabelList */
			nil,		/* by */
			reduce(8),		/* ,, reduce: LabelList */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			reduce(7),		/* $, reduce: Mutation */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S62
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			shift(7),		/* sum */
			shift(8),		/* mean */
			shift(9),		/* median */
			shift(10),		/* min */
			shift(11),		/* max */
			shift(12),		/* stddev */
			shift(14),		/* rate */
			shift(15),		/* rate_signed */
			shift(16),		/* latest */
			shift(17),		/* regularize */
			shift(18),		/* interpolate */
			shift(19),		/* percentile */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			shift(40),		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S63
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			reduce(22),		/* (, reduce: Percentile */
			nil,		/* ) */
			reduce(22),		/* by, reduce: Percentile */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
			nil,		/* [ */
			shift(79),		/* int64 */
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
			reduce(25),		/* $, reduce: Variable */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			reduce(25),		/* ,, reduce: Variable */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			shift(80),		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S67
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			shift(48),		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S68
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			reduce(28),		/* $, reduce: VarName */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			reduce(28),		/* ,, reduce: VarName */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S69
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			reduce(33),		/* ), reduce: VariableList */
			nil,		/* by */
			reduce(33),		/* ,, reduce: VariableList */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S70
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			shift(82),		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S71
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			shift(60),		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			shift(84),		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S73
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			shift(85),		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S74
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
			nil,		/* [ */
			nil,		/* int64 */
			shift(86),		/* : */
			shift(87),		/* ] */
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
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			shift(67),		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			shift(88),		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S76
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			shift(89),		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			shift(90),		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			shift(91),		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S79
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			shift(92),		/* ] */
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
			nil,		/* ) */
			nil,		/* by */
			reduce(31),		/* ,, reduce: KeyValue */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S81
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			reduce(30),		/* ,, reduce: KeyValues */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S82
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			reduce(4),		/* ), reduce: Aggregation */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			nil,		/* $ */
			nil,		/* ( */
			shift(93),		/* ) */
			nil,		/* by */
			shift(77),		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
			nil,		/* $ */
			nil,		/* ( */
			reduce(7),		/* ), reduce: Mutation */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S85
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			shift(7),		/* sum */
			shift(8),		/* mean */
			shift(9),		/* median */
			shift(10),		/* min */
			shift(11),		/* max */
			shift(12),		/* stddev */
			shift(14),		/* rate */
			shift(15),		/* rate_signed */
			shift(16),		/* latest */
			shift(17),		/* regularize */
			shift(18),		/* interpolate */
			shift(19),		/* percentile */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			shift(40),		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S86
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
			nil,		/* [ */
			shift(95),		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S87
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			reduce(25),		/* ), reduce: Variable */
			nil,		/* by */
			reduce(25),		/* ,, reduce: Variable */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S88
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			reduce(28),		/* ), reduce: VarName */
			nil,		/* by */
			reduce(28),		/* ,, reduce: VarName */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S89
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			shift(7),		/* sum */
			shift(8),		/* mean */
			shift(9),		/* median */
			shift(10),		/* min */
			shift(11),		/* max */
			shift(12),		/* stddev */
			shift(14),		/* rate */
			shift(15),		/* rate_signed */
			shift(16),		/* latest */
			shift(17),		/* regularize */
			shift(18),		/* interpolate */
			shift(19),		/* percentile */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			shift(40),		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S90
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			reduce(9),		/* ), reduce: LabelList */
			nil,		/* by */
			reduce(9),		/* ,, reduce: LabelList */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S91
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			reduce(6),		/* $, reduce: Mutation */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S92
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			reduce(24),		/* $, reduce: Variable */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			reduce(24),		/* ,, reduce: Variable */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S93
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			shift(97),		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S94
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			shift(98),		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S95
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			shift(99),		/* ] */
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S96
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			shift(100),		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S97
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			shift(7),		/* sum */
			shift(8),		/* mean */
			shift(9),		/* median */
			shift(10),		/* min */
			shift(11),		/* max */
			shift(12),		/* stddev */
			shift(14),		/* rate */
			shift(15),		/* rate_signed */
			shift(16),		/* latest */
			shift(17),		/* regularize */
			shift(18),		/* interpolate */
			shift(19),		/* percentile */
			nil,		/* [ */
			nil,		/* int64 */
			nil,		/* : */
			nil,		/* ] */
			shift(40),		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			
		},

	},
	actionRow{ // S98
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			reduce(6),		/* ), reduce: Mutation */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S99
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			reduce(24),		/* ), reduce: Variable */
			nil,		/* by */
			reduce(24),		/* ,, reduce: Variable */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S100
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			reduce(5),		/* $, reduce: Aggregation */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S101
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			shift(102),		/* ) */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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
	actionRow{ // S102
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			reduce(5),		/* ), reduce: Aggregation */
			nil,		/* by */
			nil,		/* , */
			nil,		/* varlabel */
			nil,		/* sum */
			nil,		/* mean */
			nil,		/* median */
			nil,		/* min */
			nil,		/* max */
			nil,		/* stddev */
			nil,		/* rate */
			nil,		/* rate_signed */
			nil,		/* latest */
			nil,		/* regularize */
			nil,		/* interpolate */
			nil,		/* percentile */
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

