package main

import (
    "fmt"
    "../../../model/tezos"
)

func SnapshotHeight(cycle int, snapshot int, cycle_length int,
                    snapshot_interval int) int {
    return (cycle - 7) * cycle_length + ((snapshot + 1) * snapshot_interval)
}

func BlockHashByLevel(level int) string {
    var past_head string
    head_header := tezos.Header()
    past_head = fmt.Sprintf("%s~%d", head_header.Hash, head_header.Level - level)
    past_header := tezos.HeaderAt(past_head)
    if (head_header.Level != past_header.Level) {
        fmt.Println("should not happen: tezos rpc fault, wrong level")
    }
    return past_header.Hash
}

func HashToQuery(cycle int, cycle_length int) string {
    var level_to_query int
    header := tezos.Header()
    current_level := tezos.CurrentLevelAt(header.Hash)
    blocks_ago := cycle_length * (current_level.Cycle - cycle)
    if header.Level - blocks_ago < header.Level {
        level_to_query = header.Level - blocks_ago
    } else {
	level_to_query = header.Level
    }
    return BlockHashByLevel(level_to_query)
}

func SnapshotHash(cycle int, cycle_length int, snapshot_interval int) string {
    hash := HashToQuery(cycle, cycle_length)
    cycle_info := tezos.Cycleinfo(hash, cycle)
    block_height := SnapshotHeight(cycle, cycle_info.Snapshot, cycle_length, snapshot_interval)
    return BlockHashByLevel(block_height)
}

func SnapshotLevel(cycle int, cycle_length int, snapshot_interval int) int {
    hash := HashToQuery(cycle, cycle_length)
    cycle_info := tezos.Cycleinfo(hash, cycle)
    return SnapshotHeight(cycle, cycle_info.Snapshot, cycle_length, snapshot_interval)
}

func main() {
    tezos.Initialize()

}
