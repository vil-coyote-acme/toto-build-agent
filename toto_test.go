/*
Toto-build, the stupid Go continuous build server.

Toto-build is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 3 of the License, or
(at your option) any later version.

Toto-build is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software Foundation,
Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301  USA
*/
package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"flag"
	"github.com/vil-coyote-acme/toto-build-common/testtools"
	"github.com/vil-coyote-acme/toto-build-common/message"
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"github.com/vil-coyote-acme/toto-build-common/broker"
)

func Test_Main_should_Parse_Arguments(t *testing.T) {
	// when
	main()
	// then
	assert.True(t, flag.Parsed())
	assert.Equal(t, "127.0.0.1", brokerAddr)
	assert.Equal(t, "127.0.0.1", nsqLookUpHost)
	assert.Equal(t, "4161", nsqLookUpPort)
}

func Test_Main_should_Start_An_Nsq_Service(t *testing.T) {
	//given
	if !flag.Parsed() {
		main()
	}
	b := startLookUp()
	defer b.Stop()
	// when
	startListening()
	sendMsg()
	// then
	receip, consumer := testtools.SetupListener("report")
	assert.NotNil(t, consumer)
	assert.NotNil(t, receip)
	// first get the hello from the agent
	hello := <- receip
	assert.Equal(t, "Hello", hello.Logs[0])
	// then get the build log
	buildTrace := <- receip
	assert.Contains(t, buildTrace.Logs[0], "toto-build-agent/testapp")
	close(receip)
}

func startLookUp() *broker.Broker {
	b := broker.NewBroker()
	b.StartLookUp()
	return b
}

func sendMsg() {
	// test message creation
	mess := message.ToWork{int64(1), message.TEST, "toto-build-agent/testapp"}
	body, _ := json.Marshal(mess)
	// message sending
	config := nsq.NewConfig()
	p, _ := nsq.NewProducer("127.0.0.1:4150", config)
	p.Publish("jobs", body)
}
