package paxos

import (
	"coms4113/hw5/pkg/base"
)

// Fill in the function to lead the program to a state where A2 rejects the Accept Request of P1
func ToA2RejectP1() []func(s *base.State) bool {
	p1PreparePhase := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		valid := s1.proposer.Phase == Propose
		// if valid {
		// 	fmt.Println("p1PreparePhase pass")
		// }
		return valid
	}
	p1PreparePhase_1 := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		s3 := s.Nodes()["s1"].(*Server)
		valid := s1.proposer.Phase == Propose && s1.proposer.SuccessCount == 2 && s3.n_p == 1 && s1.n_p == 1
		// if valid {
		// 	fmt.Println("p1PreparePhase_1 pass")
		// }
		return valid
	}

	p3PreparePhase := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		valid := s3.proposer.Phase == Propose && s3.n_p == 2
		// if valid {
		// 	fmt.Println("p3PreparePhase pass")
		// }
		return valid
	}
	p3PreparePhase_1 := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		valid := s3.proposer.Phase == Propose && s3.proposer.SuccessCount == 1
		// if valid {
		// 	fmt.Println("p3PreparePhase_1 pass")
		// }
		return valid
	}
	p3PreparePhase_2 := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		valid := s3.proposer.Phase == Propose && s3.proposer.SuccessCount == 2
		// if valid {
		// 	fmt.Println("p3PreparePhase_2 pass")
		// }
		return valid
	}

	p1AcceptPhase := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		valid := s1.proposer.Phase == Accept
		// if valid {
		// 	fmt.Println("p1AcceptPhase pass")
		// }
		return valid
	}
	p1AcceptPhase_1 := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		valid := s1.proposer.Phase == Accept && s1.proposer.SuccessCount == 1
		// if valid {
		// 	fmt.Println("p1AcceptPhase_1 pass")
		// }
		return valid
	}

	p3AcceptPhase := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		valid := s3.proposer.Phase == Accept
		// if valid {
		// 	fmt.Println("p3AcceptPhase pass")
		// }
		return valid
	}
	return []func(s *base.State) bool{p1PreparePhase, p1PreparePhase_1, p1AcceptPhase, p3PreparePhase, p3PreparePhase_1, p3PreparePhase_2, p1AcceptPhase, p1AcceptPhase_1, p3AcceptPhase}
}

// Fill in the function to lead the program to a state where a consensus is reached in Server 3.
func ToConsensusCase5() []func(s *base.State) bool {
	s3KnowConsensus := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		return s3.agreedValue == "v3"
	}
	return []func(s *base.State) bool{s3KnowConsensus}

}

// Fill in the function to lead the program to a state where all the Accept Requests of P1 are rejected
func NotTerminate1() []func(s *base.State) bool {

	p1 := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		s3 := s.Nodes()["s3"].(*Server)
		valid := s1.proposer.Phase == Propose && s3.proposer.Phase == "" //&& s1.proposer.SuccessCount == 3 && s1.proposer.ResponseCount == 3 && s3.proposer.Phase == ""
		// if s1.proposer.ResponseCount > 1 {
		// 	fmt.Printf("s1 response count %+v, %+v\n", s1.proposer.ResponseCount, s1.proposer.SuccessCount)
		// }
		// if valid {
		// 	fmt.Println("Phase propose s1 xx")
		// }
		return valid
	}
	a1 := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		// s2 := s.Nodes()["s2"].(*Server)
		valid := s1.proposer.Phase == Accept
		// if valid {
		// 	fmt.Println("Phase Accept s1")
		// }
		return valid
	}
	p3 := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		valid := s3.proposer.Phase == Propose && s3.proposer.SuccessCount == 2
		// if s3.proposer.ResponseCount > 1 {
		// 	fmt.Printf("s3 response count %+v, %+v\n", s3.proposer.ResponseCount, s3.proposer.SuccessCount)
		// }
		// if valid {
		// 	fmt.Println("Phase propose s3")
		// }
		return valid
	}
	p3_2 := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		// s2 := s.Nodes()["s2"].(*Server)
		valid := s3.proposer.Phase == Propose && s3.proposer.SuccessCount == 3
		// if s3.proposer.ResponseCount > 1 {
		// 	fmt.Printf("s3 response count %+v, %+v\n", s3.proposer.ResponseCount, s3.proposer.SuccessCount)
		// }
		// if valid {
		// 	fmt.Println("Phase propose s3")
		// }
		return valid
	}

	a1_1 := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		valid := s1.proposer.Phase == Accept && s1.proposer.ResponseCount == 2 && s1.proposer.SuccessCount == 0
		// if valid {
		// 	fmt.Println("Phase Accept s1 ResponseCount == 2")
		// }
		return valid
	}

	return []func(s *base.State) bool{p1, a1, p3, p3_2, a1_1}
}

// Fill in the function to lead the program to a state where all the Accept Requests of P3 are rejected
func NotTerminate2() []func(s *base.State) bool {
	p1 := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		// s3 := s.Nodes()["s3"].(*Server)
		valid := s1.proposer.Phase == Propose
		// if valid {
		// 	fmt.Println("Phase propose s1 NotTerminate2")
		// }
		return valid
	}
	a1 := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		valid := s1.proposer.Phase == Propose && s1.proposer.SuccessCount == 2
		// if valid {
		// 	fmt.Println("Phase Accept s1 NotTerminate2 SuccessCount == 2")
		// }
		return valid
	}
	a1_1 := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		valid := s1.proposer.Phase == Propose && s1.proposer.SuccessCount == 3
		// if valid {
		// 	fmt.Println("Phase Accept s1 NotTerminate2 SuccessCount == 3")
		// }
		return valid
	}
	a3 := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		valid := s3.proposer.Phase == Accept
		// if valid {
		// 	fmt.Println("Phase Accept s3")
		// }
		return valid
	}
	a3_1 := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		valid := s3.proposer.Phase == Accept && s3.proposer.ResponseCount == 1 && s3.proposer.SuccessCount == 0
		// if valid {
		// 	fmt.Println("Phase Accept s3 ResponseCount == 1")
		// }
		return valid
	}
	a3_2 := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		valid := s3.proposer.Phase == Accept && s3.proposer.ResponseCount == 2 && s3.proposer.SuccessCount == 0
		// if valid {
		// 	fmt.Println("Phase Accept s3 ResponseCount == 2")
		// }
		return valid
	}

	return []func(s *base.State) bool{p1, a1, a1_1, a3, a3_1, a3_2}
}

// Fill in the function to lead the program to a state where all the Accept Requests of P1 are rejected again.
func NotTerminate3() []func(s *base.State) bool {
	p1 := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		// s3 := s.Nodes()["s3"].(*Server)
		valid := s1.proposer.Phase == Propose
		// if s1.proposer.ResponseCount > 1 {
		// 	fmt.Printf("s1 response count %+v, %+v\n", s1.proposer.ResponseCount, s1.proposer.SuccessCount)
		// }
		// if valid {
		// 	fmt.Println("Phase propose s1 xx")
		// }
		return valid
	}
	a1 := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		// s2 := s.Nodes()["s2"].(*Server)
		valid := s1.proposer.Phase == Accept
		// if valid {
		// 	fmt.Println("Phase Accept s1")
		// }
		return valid
	}
	p3 := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		valid := s3.proposer.Phase == Propose && s3.proposer.SuccessCount == 2
		// if s3.proposer.ResponseCount > 1 {
		// 	fmt.Printf("s3 response count %+v, %+v\n", s3.proposer.ResponseCount, s3.proposer.SuccessCount)
		// }
		// if valid {
		// 	fmt.Println("Phase propose s3")
		// }
		return valid
	}
	p3_2 := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		// s2 := s.Nodes()["s2"].(*Server)
		valid := s3.proposer.Phase == Propose && s3.proposer.SuccessCount == 3
		// if s3.proposer.ResponseCount > 1 {
		// 	fmt.Printf("s3 response count %+v, %+v\n", s3.proposer.ResponseCount, s3.proposer.SuccessCount)
		// }
		// if valid {
		// 	fmt.Println("Phase propose s3")
		// }
		return valid
	}

	a1_1 := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		valid := s1.proposer.Phase == Accept && s1.proposer.ResponseCount == 2 && s1.proposer.SuccessCount == 0
		// if valid {
		// 	fmt.Println("Phase Accept s1 ResponseCount == 2")
		// }
		return valid
	}

	return []func(s *base.State) bool{p1, a1, p3, p3_2, a1_1}
}

// Fill in the function to lead the program to make P1 propose first, then P3 proposes, but P1 get rejects in
// Accept phase
func concurrentProposer1() []func(s *base.State) bool {
	p1PreparePhase := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		valid := s1.proposer.Phase == Propose && s1.proposer.N == 1
		// if valid {
		// 	fmt.Println("p1PreparePhase pass")
		// }
		return valid
	}
	p1PreparePhase_1 := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		valid := s1.proposer.Phase == Propose && s1.proposer.SuccessCount == 2 && s1.proposer.N == 1
		// if valid {
		// 	fmt.Println("p1PreparePhase_1 pass")
		// }
		return valid
	}
	p1PreparePhase_2 := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		valid := s1.proposer.Phase == Propose && s1.proposer.SuccessCount == 3 && s1.proposer.N == 1
		// if valid {
		// 	fmt.Println("p1PreparePhase_2 pass")
		// }
		return valid
	}
	p3PreparePhase := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		s1 := s.Nodes()["s1"].(*Server)
		// s1 := s.Nodes()["s1"].(*Server)
		valid := s3.proposer.Phase == Propose && s3.proposer.N == 2 && s1.proposer.SuccessCount == 3
		// if valid {
		// 	fmt.Println("p3PreparePhase pass")
		// 	fmt.Printf("s3 pass %+v \n", s3)
		// 	fmt.Printf("s1 pass %+v\n", s1)
		// }
		return valid
	}
	p3PreparePhase_1 := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		s1 := s.Nodes()["s1"].(*Server)
		valid := s3.proposer.Phase == Propose && s3.proposer.SuccessCount == 2 && s1.proposer.Phase == Propose
		// if valid {
		// 	fmt.Println("p3PreparePhase_1 pass")

		// }
		return valid
	}
	p3PreparePhase_2 := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		s1 := s.Nodes()["s1"].(*Server)
		valid := s3.proposer.Phase == Propose && s3.proposer.SuccessCount == 3 && s1.proposer.Phase == Propose
		// if valid {
		// 	fmt.Println("p3PreparePhase_2 pass")
		// 	fmt.Printf("s3 pass %+v \n", s3)
		// 	fmt.Printf("s1 pass %+v\n", s1)

		// }
		return valid
	}

	p1Acceptphase := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		s3 := s.Nodes()["s3"].(*Server)

		valid := s1.proposer.Phase == Accept && s3.proposer.Phase == Propose && s3.proposer.N == 2 && s1.proposer.N == 1
		// if valid {
		// 	fmt.Println("p1Acceptphase pass")
		// 	fmt.Printf("s3 pass %+v \n", s3)
		// 	fmt.Printf("s1 pass %+v\n", s1)

		// }
		return valid
	}
	p1Acceptphase_1 := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		s3 := s.Nodes()["s3"].(*Server)

		valid := s1.proposer.Phase == Accept && s3.proposer.Phase == Propose && s1.proposer.SuccessCount == 0 && s1.proposer.ResponseCount == 2
		// if valid {
		// 	fmt.Println("p1Acceptphase_1 pass")
		// 	fmt.Printf("s3 pass %+v \n", s3)
		// 	fmt.Printf("s1 pass %+v\n", s1)

		// }
		return valid
	}
	p1Acceptphase_2 := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		s3 := s.Nodes()["s3"].(*Server)

		valid := s1.proposer.Phase == Accept && s3.proposer.Phase == Propose && s1.proposer.SuccessCount == 0 && s1.proposer.ResponseCount == 3
		// if valid {
		// 	fmt.Println("p1Acceptphase_2 pass")
		// 	fmt.Printf("s3 pass %+v \n", s3)
		// 	fmt.Printf("s1 pass %+v\n", s1)

		// }
		return valid
	}

	return []func(s *base.State) bool{p1PreparePhase, p1PreparePhase_1, p1PreparePhase_2, p3PreparePhase, p3PreparePhase_1, p3PreparePhase_2, p1Acceptphase, p1Acceptphase_1, p1Acceptphase_2}

}

// Fill in the function to lead the program continue  P3's proposal  and reaches consensus at the value of "v3".
func concurrentProposer2() []func(s *base.State) bool {
	p3Acceptphase := func(s *base.State) bool {
		// s1 := s.Nodes()["s1"].(*Server)
		s3 := s.Nodes()["s3"].(*Server)

		valid := s3.proposer.Phase == Accept && s3.proposer.SuccessCount == 2
		// if valid {
		// 	fmt.Println("p1Acceptphase pass")
		// 	fmt.Printf("s3 pass %+v \n", s3)
		// 	fmt.Printf("s1 pass %+v\n", s1)

		// }
		return valid
	}
	// p3Acceptphase := func(s *base.State) bool {
	// 	s1 := s.Nodes()["s1"].(*Server)
	// 	s3 := s.Nodes()["s3"].(*Server)

	// 	valid := s3.proposer.Phase == Accept
	// 	if valid {
	// 		fmt.Println("p1Acceptphase pass")
	// 		fmt.Printf("s3 pass %+v \n", s3)
	// 		fmt.Printf("s1 pass %+v\n", s1)

	// 	}
	// 	return valid
	// }
	return []func(s *base.State) bool{p3Acceptphase}

}
