package miner

import (
	"log"
	"strconv"
	"strings"
)

const (
	JOINOK     = "JOINOK"
	PASSFAILED = "PASSFAILED"
	PAYMENTOK  = "PAYMENTOK"
	PONG       = "PONG"
	POOLSTEPS  = "POOLSTEPS"
	STATUS     = "STATUS"
	STEPOK     = "STEPOK"
	STEPFAIL   = "STEPFAIL"
)

func Parse(comms *Comms, poolIp string, wallet string, block int, resp string) {
	if resp == "" {
		log.Println("Got an empty response")
		return
	}
	r := strings.Split(resp, " ")

	switch r[0] {
	case JOINOK:
		comms.PoolAddr <- r[1]
		comms.MinerSeed <- r[2]
		poolData(comms, r, 2)
		comms.Joined <- struct{}{}
	case PASSFAILED:
		log.Println("Incorrect pool password")
	case PAYMENTOK:
		LogPaymentResp(r, poolIp)
	case PONG:
		comms.Pong <- struct{}{}
	case POOLSTEPS:
		poolData(comms, r, 0)
	case STEPOK:
		shares, err := strconv.Atoi(r[1])
		if err != nil {
			log.Printf("Had trouble parsing the shares from STEPOK message: %v\n", err)
		}
		comms.StepSolved <- shares
	case STEPFAIL:
		comms.StepFailed <- 1
	case STATUS:
		comms.PoolStatus <- NewPoolStatus(r[1:])
	default:
		log.Printf("Uknown response code: %s\n", r[0])
	}
}

func poolData(comms *Comms, resp []string, offset int) {
	block, err := strconv.Atoi(resp[2+offset])
	if err != nil {
		log.Printf("Error converting target block: %s\n", resp[2+offset])
	} else {
		comms.Block <- block
	}

	comms.TargetString <- resp[3+offset]

	targetChars, err := strconv.Atoi(resp[4+offset])
	if err != nil {
		log.Printf("Error converting target chars: %s\n", resp[4+offset])
	} else {
		comms.TargetChars <- targetChars
	}

	step, err := strconv.Atoi(resp[5+offset])
	if err != nil {
		log.Printf("Error converting target chars: %s\n", resp[5+offset])
	} else {
		comms.Step <- step
	}

	diff, err := strconv.Atoi(resp[6+offset])
	if err != nil {
		log.Printf("Error converting target chars: %s\n", resp[6+offset])
	} else {
		comms.Diff <- diff
	}

	comms.Balance <- resp[7+offset]

	blocksTillPayment, err := strconv.Atoi(resp[8+offset])
	if err != nil {
		log.Printf("Error converting target chars: %s\n", resp[8+offset])
	} else {
		comms.BlocksTillPayment <- blocksTillPayment
	}

	comms.PoolHashRate <- resp[9+offset] + "000"

	poolDepth, err := strconv.Atoi(resp[10+offset])
	if err != nil {
		log.Printf("Error converting target chars: %s\n", resp[10+offset])
	} else {
		comms.PoolDepth <- poolDepth
	}
}
