package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	buflen = 512      //Max length of buffer
	port   = ":25045" // = "JBP" in base 36!
)

// Watchvar holds data about a watched variable
type Watchvar struct {
	ref    *int
	oldval int
	name   string
}

// Controlvar holds data about a controlled variable
type Controlvar struct {
	ref    *int
	oldval int
	maxval int
	minval int
}

var monitor struct {
	watchvars   []Watchvar
	controlvars map[string]Controlvar
	remoteAddr  *net.UDPAddr
	localAddr   *net.UDPAddr
	running     bool
}

// MonitorInit initialises the monitor function
func MonitorInit() {
	var err error
	monitor.controlvars = make(map[string]Controlvar)
	monitor.localAddr, err = net.ResolveUDPAddr("udp", port)
	if err != nil {
		log.Fatal(err)
	}
	monitor.running = true

	go listener()
	go watcher()
}

// Watch starts watching a variable
func Watch(watchvar *int, name string) {
	monitor.watchvars = append(monitor.watchvars, Watchvar{ref: watchvar, oldval: *watchvar, name: name})
}

// Control registers a control variable for remote control
func Control(controlvar *int, name string, maxval int, minval int) {
	monitor.controlvars[name] = Controlvar{ref: controlvar, oldval: *controlvar, maxval: maxval, minval: minval}
}

// JSONControl struct used for communication
type JSONControl struct {
	min int
	max int
	val int
}

func listener() {
	buf := make([]byte, buflen)

	/* Listen at selected port */
	ServerConn, err := net.ListenUDP("udp", monitor.localAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer ServerConn.Close()

	for monitor.running {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		if err != nil {
			log.Println("Error: ", err)
		} else {
			monitor.remoteAddr = addr
		}
		log.Println("Received ", string(buf[0:n]), " from ", addr)

		var jdat map[string]interface{}
		if err = json.Unmarshal(buf[0:n], &jdat); err != nil {
			panic(err)
		}
		fmt.Println(jdat)

		/*
			// Look for any "set"s of control variables
			if (j.count("set") > 0) {
					for (json::iterator it = j["set"].begin(); it != j["set"].end(); ++it) {
							if (controlvars.count(it.key()) > 0) {
									*(controlvars[it.key()].ref) = it.value();
							}
					}
			}
		*/
		// Send back the current set of control variables
		jout := make(map[string]JSONControl)
		for name, c := range monitor.controlvars {
			jout[name] = JSONControl{min: c.minval, max: c.maxval, val: 0}
		}
		jdat["control"] = jout

		// Convert the data into JSON
		joutbytes, err := json.Marshal(jdat)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(joutbytes))

		//Reply to the client
		if _, err := ServerConn.WriteToUDP(joutbytes, monitor.remoteAddr); err != nil {
			log.Println("Error: ", err)
		}
	}
}

const updateinterval = 200

func watcher() {

	//loopcount := 0

	for monitor.running {
		time.Sleep(2 * time.Millisecond)

	}

	/*
	   	json j;
	      bool dataToSend;
	      unsigned int loopcount = 0;
	      int watchSocket_fd = *((int*) param);

	      while (running)
	      {
	          usleep(2000);  // loop once per 2ms
	          loopcount++;
	          j.clear();
	          dataToSend = false;
	          j["millis"] = millis();
	          int len = watchvars.size();
	          for (int i = 0; i < len; i++) {
	              if ((loopcount % UPDATEINTERVAL) == 0 || watchvars[i].oldval != *watchvars[i].ref) {
	                  // Add the watched variable to the report
	                  j["variables"][watchvars[i].name] = *watchvars[i].ref;
	                  watchvars[i].oldval = *watchvars[i].ref;
	                  dataToSend = true;
	              }
	          }
	          if (dataToSend) {
	              std::string dataout(j.dump());
	              if (socket_other.sin_addr.s_addr != 0) {
	                  if (sendto(watchSocket_fd, dataout.c_str(), dataout.length(), 0, (struct sockaddr *) &socket_other, slen) == -1) {
	                      die("sendto()");
	                  }
	              }
	          }
	      }
	      close(watchSocket_fd);
	      return 0;
	*/

}
