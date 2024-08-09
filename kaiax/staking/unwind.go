package staking

func (s *StakingModule) Unwind(num uint64) error {
	if !s.isKaia(num) && (num%s.stakingInterval) == 0 {
		DeleteStakingInfo(s.ChainKv, num)
	}
	return nil
}
