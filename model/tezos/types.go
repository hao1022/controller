package tezos

type BlockHeaderType struct{
    Protocol         string   `json:"protocol"`
    ChainId          string   `json:"chain_id"`
    Hash             string   `json:"hash"`
    Level            int      `json:"level"`
    Proto            int      `json:"proto"`
    Predecessor      string   `json:"predecessor"`
    Timestamp        string   `json:"timestamp"`
    ValidationPass   int      `json:"validation_pass"`
    OperationsHash   string   `json:"operations_hash"`
    Fitness          []string `json:"fitness"`
    Context          string   `json:"context"`
    Priority         int      `json:"priority"`
    Nonce            string   `json:"proof_of_work_nonce"`
    Signature        string   `json:"signature"`
}

type BalanceUpdateType struct{
    Kind     string   `json:"kind"`
    Category string   `json:"category"`
    Contract string   `json:"contract"`
    Delegate string   `json:"delegate"`
    Level    int      `json:"level"`
    Change   string   `json:"change"`
}

type BlockMetadataType struct{
    Protocol       string                 `json:"protocol"`
    Baker          string                 `json:"baker"`
    BalanceUpdates []BalanceUpdateType    `json:"balance_updates"`
}

type EndorsingRightType struct{
    Level         int      `json:"level"`
    Delegate      string   `json:"delegate"`
    Slots         []int    `json:"slots"`
    EstimatedTime int      `json:"estimated_time"`
}

type BakingRightType struct{
    Level         int      `json:"level"`
    Delegate      string   `json:"delegate"`
    Priority      int      `json:"priority"`
    EstimatedTime int      `json:"estimated_time"`
}

type CurrentLevelType struct{
    Level                int   `json:"level"`
    LevelPosition        int   `json:"level_position"`
    Cycle                int   `json:"cycle"`
    CyclePosition        int   `json:"cycle_position"`
    VotingPeriod         int   `json:"voting_period"`
    VotingPeriodPosition int   `json:"voting_period_position"`
    ExpectedCommitment   bool  `json:"expected_commitment"`
}

type OperationResultType struct{
    Status      string  `json:"status"`
    ConsumedGas string  `json:"consumed_gas"`
}

type OperationMetadataType struct{
    Result OperationResultType  `json:"operation_result"`
}

type OperationContentType struct{
    Kind        string   `json:"kind"`
    Source      string   `json:"source"`
    Fee         string   `json:"fee"`
    Level       int      `json:"level"`
    Contract    string   `json:"contract"`
    Change      string   `json:"change"`
    Metadata    OperationMetadataType `json:"metadata"`
    Destination string   `json:"destination"`
    Amount      string   `json:"amount"`
}

type OperationType struct{
    Protocol string   `json:"protocol"`
    Hash     string   `json:"hash"`
    Branch   string   `json:"branch"`
    Contents []OperationContentType  `json:"contents"`
}

type FrozenBalanceByCycleType struct{
    Cycle   int      `json:"cycle"`
    Deposit string   `json:"deposit"`
    Fees    string   `json:"fees"`
    Rewards string   `json:"rewards"`
}

type CycleInfoType struct{
    LastRoll   []string `json:"last_roll"`
    Nonces     []string `json:"nonces"`
    RandomSeed string   `json:"random_seed"`
    Snapshot   int      `json:"roll_snapshot"`
}
