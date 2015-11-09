
package lexer



/*
Let s be the current state
Let r be the current input rune
transitionTable[s](r) returns the next state.
*/
type TransitionTable [NumStates] func(rune) int

var TransTab = TransitionTable{
	
		// S0
		func(r rune) int {
			switch {
			case r == 9 : // ['\t','\t']
				return 1
			case r == 10 : // ['\n','\n']
				return 1
			case r == 13 : // ['\r','\r']
				return 1
			case r == 32 : // [' ',' ']
				return 1
			case r == 40 : // ['(','(']
				return 2
			case r == 41 : // [')',')']
				return 3
			case r == 44 : // [',',',']
				return 4
			case r == 47 : // ['/','/']
				return 5
			case r == 48 : // ['0','0']
				return 6
			case 49 <= r && r <= 57 : // ['1','9']
				return 7
			case r == 61 : // ['=','=']
				return 8
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 97 : // ['a','a']
				return 9
			case r == 98 : // ['b','b']
				return 10
			case r == 99 : // ['c','c']
				return 9
			case r == 100 : // ['d','d']
				return 11
			case 101 <= r && r <= 107 : // ['e','k']
				return 9
			case r == 108 : // ['l','l']
				return 12
			case r == 109 : // ['m','m']
				return 13
			case 110 <= r && r <= 111 : // ['n','o']
				return 9
			case r == 112 : // ['p','p']
				return 14
			case r == 113 : // ['q','q']
				return 9
			case r == 114 : // ['r','r']
				return 15
			case r == 115 : // ['s','s']
				return 16
			case 116 <= r && r <= 122 : // ['t','z']
				return 9
			case r == 123 : // ['{','{']
				return 17
			case r == 125 : // ['}','}']
				return 18
			
			
			
			}
			return NoState
			
		},
	
		// S1
		func(r rune) int {
			switch {
			
			
			
			}
			return NoState
			
		},
	
		// S2
		func(r rune) int {
			switch {
			
			
			
			}
			return NoState
			
		},
	
		// S3
		func(r rune) int {
			switch {
			
			
			
			}
			return NoState
			
		},
	
		// S4
		func(r rune) int {
			switch {
			
			
			
			}
			return NoState
			
		},
	
		// S5
		func(r rune) int {
			switch {
			case r == 42 : // ['*','*']
				return 19
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 21
			case 48 <= r && r <= 57 : // ['0','9']
				return 22
			case 65 <= r && r <= 90 : // ['A','Z']
				return 23
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 122 : // ['a','z']
				return 23
			
			
			
			}
			return NoState
			
		},
	
		// S6
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 122 : // ['a','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S7
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 7
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 122 : // ['a','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S8
		func(r rune) int {
			switch {
			
			
			
			}
			return NoState
			
		},
	
		// S9
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 122 : // ['a','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S10
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 120 : // ['a','x']
				return 9
			case r == 121 : // ['y','y']
				return 24
			case r == 122 : // ['z','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S11
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 100 : // ['a','d']
				return 9
			case r == 101 : // ['e','e']
				return 25
			case 102 <= r && r <= 122 : // ['f','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S12
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case r == 97 : // ['a','a']
				return 26
			case 98 <= r && r <= 122 : // ['b','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S13
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case r == 97 : // ['a','a']
				return 27
			case 98 <= r && r <= 100 : // ['b','d']
				return 9
			case r == 101 : // ['e','e']
				return 28
			case 102 <= r && r <= 104 : // ['f','h']
				return 9
			case r == 105 : // ['i','i']
				return 29
			case 106 <= r && r <= 122 : // ['j','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S14
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 100 : // ['a','d']
				return 9
			case r == 101 : // ['e','e']
				return 30
			case 102 <= r && r <= 122 : // ['f','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S15
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case r == 97 : // ['a','a']
				return 31
			case 98 <= r && r <= 122 : // ['b','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S16
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 115 : // ['a','s']
				return 9
			case r == 116 : // ['t','t']
				return 32
			case r == 117 : // ['u','u']
				return 33
			case 118 <= r && r <= 122 : // ['v','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S17
		func(r rune) int {
			switch {
			case r == 125 : // ['}','}']
				return 34
			
			
			
			}
			return NoState
			
		},
	
		// S18
		func(r rune) int {
			switch {
			
			
			
			}
			return NoState
			
		},
	
		// S19
		func(r rune) int {
			switch {
			case r == 42 : // ['*','*']
				return 35
			
			
			default:
				return 19
			}
			
		},
	
		// S20
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 122 : // ['a','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S21
		func(r rune) int {
			switch {
			case r == 10 : // ['\n','\n']
				return 36
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 122 : // ['a','z']
				return 9
			
			
			default:
				return 37
			}
			
		},
	
		// S22
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 38
			case r == 47 : // ['/','/']
				return 38
			case 48 <= r && r <= 57 : // ['0','9']
				return 22
			case 65 <= r && r <= 90 : // ['A','Z']
				return 23
			case r == 95 : // ['_','_']
				return 38
			case 97 <= r && r <= 122 : // ['a','z']
				return 23
			
			
			
			}
			return NoState
			
		},
	
		// S23
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 38
			case r == 47 : // ['/','/']
				return 38
			case 48 <= r && r <= 57 : // ['0','9']
				return 22
			case 65 <= r && r <= 90 : // ['A','Z']
				return 23
			case r == 95 : // ['_','_']
				return 38
			case 97 <= r && r <= 122 : // ['a','z']
				return 23
			
			
			
			}
			return NoState
			
		},
	
		// S24
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 122 : // ['a','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S25
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 107 : // ['a','k']
				return 9
			case r == 108 : // ['l','l']
				return 39
			case 109 <= r && r <= 122 : // ['m','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S26
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 115 : // ['a','s']
				return 9
			case r == 116 : // ['t','t']
				return 40
			case 117 <= r && r <= 122 : // ['u','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S27
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 119 : // ['a','w']
				return 9
			case r == 120 : // ['x','x']
				return 41
			case 121 <= r && r <= 122 : // ['y','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S28
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case r == 97 : // ['a','a']
				return 42
			case 98 <= r && r <= 99 : // ['b','c']
				return 9
			case r == 100 : // ['d','d']
				return 43
			case 101 <= r && r <= 122 : // ['e','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S29
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 109 : // ['a','m']
				return 9
			case r == 110 : // ['n','n']
				return 44
			case 111 <= r && r <= 122 : // ['o','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S30
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 113 : // ['a','q']
				return 9
			case r == 114 : // ['r','r']
				return 45
			case 115 <= r && r <= 122 : // ['s','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S31
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 115 : // ['a','s']
				return 9
			case r == 116 : // ['t','t']
				return 46
			case 117 <= r && r <= 122 : // ['u','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S32
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 99 : // ['a','c']
				return 9
			case r == 100 : // ['d','d']
				return 47
			case 101 <= r && r <= 122 : // ['e','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S33
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 108 : // ['a','l']
				return 9
			case r == 109 : // ['m','m']
				return 48
			case 110 <= r && r <= 122 : // ['n','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S34
		func(r rune) int {
			switch {
			
			
			
			}
			return NoState
			
		},
	
		// S35
		func(r rune) int {
			switch {
			case r == 42 : // ['*','*']
				return 35
			case r == 47 : // ['/','/']
				return 49
			
			
			default:
				return 19
			}
			
		},
	
		// S36
		func(r rune) int {
			switch {
			
			
			
			}
			return NoState
			
		},
	
		// S37
		func(r rune) int {
			switch {
			case r == 10 : // ['\n','\n']
				return 36
			
			
			default:
				return 37
			}
			
		},
	
		// S38
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 38
			case r == 47 : // ['/','/']
				return 38
			case 48 <= r && r <= 57 : // ['0','9']
				return 22
			case 65 <= r && r <= 90 : // ['A','Z']
				return 23
			case r == 95 : // ['_','_']
				return 38
			case 97 <= r && r <= 122 : // ['a','z']
				return 23
			
			
			
			}
			return NoState
			
		},
	
		// S39
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 115 : // ['a','s']
				return 9
			case r == 116 : // ['t','t']
				return 50
			case 117 <= r && r <= 122 : // ['u','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S40
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 100 : // ['a','d']
				return 9
			case r == 101 : // ['e','e']
				return 51
			case 102 <= r && r <= 122 : // ['f','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S41
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 122 : // ['a','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S42
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 109 : // ['a','m']
				return 9
			case r == 110 : // ['n','n']
				return 52
			case 111 <= r && r <= 122 : // ['o','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S43
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 104 : // ['a','h']
				return 9
			case r == 105 : // ['i','i']
				return 53
			case 106 <= r && r <= 122 : // ['j','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S44
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 122 : // ['a','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S45
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 98 : // ['a','b']
				return 9
			case r == 99 : // ['c','c']
				return 54
			case 100 <= r && r <= 122 : // ['d','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S46
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 100 : // ['a','d']
				return 9
			case r == 101 : // ['e','e']
				return 55
			case 102 <= r && r <= 122 : // ['f','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S47
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 99 : // ['a','c']
				return 9
			case r == 100 : // ['d','d']
				return 56
			case 101 <= r && r <= 122 : // ['e','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S48
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 122 : // ['a','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S49
		func(r rune) int {
			switch {
			
			
			
			}
			return NoState
			
		},
	
		// S50
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case r == 97 : // ['a','a']
				return 57
			case 98 <= r && r <= 122 : // ['b','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S51
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 114 : // ['a','r']
				return 9
			case r == 115 : // ['s','s']
				return 58
			case 116 <= r && r <= 122 : // ['t','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S52
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 122 : // ['a','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S53
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case r == 97 : // ['a','a']
				return 59
			case 98 <= r && r <= 122 : // ['b','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S54
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 100 : // ['a','d']
				return 9
			case r == 101 : // ['e','e']
				return 60
			case 102 <= r && r <= 122 : // ['f','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S55
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 61
			case 97 <= r && r <= 122 : // ['a','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S56
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 100 : // ['a','d']
				return 9
			case r == 101 : // ['e','e']
				return 62
			case 102 <= r && r <= 122 : // ['f','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S57
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 122 : // ['a','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S58
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 115 : // ['a','s']
				return 9
			case r == 116 : // ['t','t']
				return 63
			case 117 <= r && r <= 122 : // ['u','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S59
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 109 : // ['a','m']
				return 9
			case r == 110 : // ['n','n']
				return 64
			case 111 <= r && r <= 122 : // ['o','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S60
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 109 : // ['a','m']
				return 9
			case r == 110 : // ['n','n']
				return 65
			case 111 <= r && r <= 122 : // ['o','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S61
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 114 : // ['a','r']
				return 9
			case r == 115 : // ['s','s']
				return 66
			case 116 <= r && r <= 122 : // ['t','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S62
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 117 : // ['a','u']
				return 9
			case r == 118 : // ['v','v']
				return 67
			case 119 <= r && r <= 122 : // ['w','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S63
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 122 : // ['a','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S64
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 122 : // ['a','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S65
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 115 : // ['a','s']
				return 9
			case r == 116 : // ['t','t']
				return 68
			case 117 <= r && r <= 122 : // ['u','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S66
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 104 : // ['a','h']
				return 9
			case r == 105 : // ['i','i']
				return 69
			case 106 <= r && r <= 122 : // ['j','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S67
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 122 : // ['a','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S68
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 104 : // ['a','h']
				return 9
			case r == 105 : // ['i','i']
				return 70
			case 106 <= r && r <= 122 : // ['j','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S69
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 102 : // ['a','f']
				return 9
			case r == 103 : // ['g','g']
				return 71
			case 104 <= r && r <= 122 : // ['h','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S70
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 107 : // ['a','k']
				return 9
			case r == 108 : // ['l','l']
				return 72
			case 109 <= r && r <= 122 : // ['m','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S71
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 109 : // ['a','m']
				return 9
			case r == 110 : // ['n','n']
				return 73
			case 111 <= r && r <= 122 : // ['o','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S72
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 100 : // ['a','d']
				return 9
			case r == 101 : // ['e','e']
				return 74
			case 102 <= r && r <= 122 : // ['f','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S73
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 100 : // ['a','d']
				return 9
			case r == 101 : // ['e','e']
				return 75
			case 102 <= r && r <= 122 : // ['f','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S74
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 122 : // ['a','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S75
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 99 : // ['a','c']
				return 9
			case r == 100 : // ['d','d']
				return 76
			case 101 <= r && r <= 122 : // ['e','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
		// S76
		func(r rune) int {
			switch {
			case r == 45 : // ['-','-']
				return 20
			case r == 47 : // ['/','/']
				return 20
			case 48 <= r && r <= 57 : // ['0','9']
				return 6
			case 65 <= r && r <= 90 : // ['A','Z']
				return 9
			case r == 95 : // ['_','_']
				return 20
			case 97 <= r && r <= 122 : // ['a','z']
				return 9
			
			
			
			}
			return NoState
			
		},
	
}
