/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceSDK contributors
 SPDX-License-Identifier: Apache-2.0
*/

package utils

import (
	"crypto/rand"
	"math/big"
	"strconv"
	"sync/atomic"
	"time"
)

// TicketIDAtomicCounter is a counter that adds when each time a function is called
var TicketIDAtomicCounter int64

// TicketIDGenerator is a random unique ticket ID generator, output is a string
func TicketIDGenerator(prefix string) Ticket {
	now := 10000000 + int(time.Now().UnixNano())%89999999

	atomic.AddInt64(&TicketIDAtomicCounter, 1)
	atomicCounter := 10000 + int(TicketIDAtomicCounter)%89999

	random, Error := rand.Int(rand.Reader, big.NewInt(89999))
	if Error != nil {
		panic(Error)
	}

	randomNumber := 10000 + random.Int64()
	trulyRandNumber := prefix + strconv.Itoa(atomicCounter) + strconv.Itoa(now) + strconv.FormatInt(randomNumber, 10)
	ticket := Ticket(trulyRandNumber)

	return ticket
}
