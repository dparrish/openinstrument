
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
			shift(19),		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S5
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			shift(20),		/* ( */
			nil,		/* ) */
			shift(21),		/* by */
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S6
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			shift(22),		/* ( */
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S8
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
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
			nil,		/* varname */
			shift(24),		/* {} */
			shift(25),		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S20
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
			shift(30),		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S21
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			shift(31),		/* ( */
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S22
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			shift(33),		/* varlabel */
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			shift(35),		/* int64 */
			
		},

	},
	actionRow{ // S24
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S26
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			reduce(29),		/* ), reduce: VariableList */
			nil,		/* by */
			reduce(29),		/* ,, reduce: VariableList */
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S27
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			shift(39),		/* ) */
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S28
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			shift(40),		/* ) */
			nil,		/* by */
			shift(41),		/* , */
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S29
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			shift(42),		/* ( */
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S30
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
			nil,		/* varname */
			shift(43),		/* {} */
			shift(44),		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S31
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			nil,		/* , */
			shift(46),		/* varlabel */
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S32
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			shift(47),		/* , */
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S33
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			reduce(31),		/* ,, reduce: Duration */
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S34
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			shift(48),		/* ) */
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S35
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			reduce(32),		/* ), reduce: Int64 */
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			shift(49),		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S37
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			shift(50),		/* , */
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			shift(51),		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S38
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			reduce(26),		/* ,, reduce: KeyValues */
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			reduce(26),		/* }, reduce: KeyValues */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S39
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S40
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
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
			shift(30),		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
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
			nil,		/* , */
			shift(33),		/* varlabel */
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S43
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S45
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			shift(55),		/* ) */
			nil,		/* by */
			shift(56),		/* , */
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S46
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
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
			shift(30),		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S48
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
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
			shift(58),		/* varlabel */
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S51
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S52
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			reduce(30),		/* ), reduce: VariableList */
			nil,		/* by */
			reduce(30),		/* ,, reduce: VariableList */
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
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
			shift(60),		/* , */
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
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
			shift(50),		/* , */
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			shift(61),		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S55
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			shift(62),		/* ( */
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
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
			shift(63),		/* varlabel */
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S57
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			shift(64),		/* ) */
			nil,		/* by */
			shift(41),		/* , */
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
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
			reduce(28),		/* ,, reduce: KeyValue */
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			reduce(28),		/* }, reduce: KeyValue */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S59
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			nil,		/* ) */
			nil,		/* by */
			reduce(27),		/* ,, reduce: KeyValues */
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			reduce(27),		/* }, reduce: KeyValues */
			nil,		/* = */
			nil,		/* int64 */
			
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
			shift(30),		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S61
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
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
			shift(30),		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S63
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S64
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S65
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			shift(68),		/* ) */
			nil,		/* by */
			shift(41),		/* , */
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S66
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			shift(69),		/* ) */
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S67
				canRecover: false,
		actions: [numSymbols]action{
			nil,		/* INVALID */
			nil,		/* $ */
			nil,		/* ( */
			shift(70),		/* ) */
			nil,		/* by */
			shift(41),		/* , */
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S68
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S69
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	actionRow{ // S70
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
			nil,		/* varname */
			nil,		/* {} */
			nil,		/* { */
			nil,		/* } */
			nil,		/* = */
			nil,		/* int64 */
			
		},

	},
	
}

