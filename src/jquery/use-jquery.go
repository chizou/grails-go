package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const numColorfulDisplayDots = 100 //16
var q [numColorfulDisplayDots] *testVector
var qhealth = newHealthArray()
var containers []string

func home(w http.ResponseWriter, r *http.Request) {
	containers = r.URL.Query()["containers"]

	html := `<head>	
 <script src="http://cdnjs.cloudflare.com/ajax/libs/raphael/2.1.0/raphael-min.js">
 </script>

 <script src='http://ajax.googleapis.com/ajax/libs/jquery/1.11.2/jquery.min.js'>
 </script>

 </head>
 <body>
  <script>
 var si1

 function getMessages2() {
  $.post('status', {}, function(r) {
   var foox = JSON.parse(r)
   r = foox.dotStatus
   t = foox.dotDuration
   v = foox.squareHealth
   paper.text(800, 400, v);
   x = foox.dotColor
   var THROWAWAY_BR_CHARS = 0
   var ROWS_PER_SQUARE = Math.sqrt(r.length) //10
   var COLS_PER_SQUARE = ROWS_PER_SQUARE
   var ORIGINX = 40
   var ORIGINY = 100
   var STOPLIGHT_YELLOW = '#FAD201'
   var STOPLIGHT_GREEN = '#27E833'
   var MEDIUM_GREEN =    '#27C833'
   var DARK_GREEN =      '#27A833'
   var STOPLIGHT_BLUE = '#CCE5FF'
   var MEDIUM_BLUE =    '#66B2FF'
   var DARK_BLUE =      '#0080FF'
   for (row = 0; row < ROWS_PER_SQUARE; row++) {
    for (col = 0; col < COLS_PER_SQUARE; col++) {
       var stroke
       var fill
       var filltext = ''
       var substring_start = THROWAWAY_BR_CHARS+row*ROWS_PER_SQUARE+col
       //console.log(substring_start)
       var rchar = r.substring(substring_start, substring_start+1)
       if (rchar == '0') {
           stroke = "white"
           fill = "white"
       }
       if (rchar == '1') {
           stroke = STOPLIGHT_YELLOW
           fill = "white"
       }
       if (rchar == '2') {
           stroke = STOPLIGHT_YELLOW
           fill = STOPLIGHT_YELLOW
       }
       if (rchar == '4') {
           fill = 'pink'
           filltext = 'error'
           var tchar = t.substring(substring_start, substring_start+1)
           var xchar = x.substring(substring_start, substring_start+1)
           if ((tchar == '0') ||
               (tchar == '1')) {
               if (xchar == 'G') {
                   fill = DARK_GREEN
               }
               if (xchar == 'B') {
                   fill = DARK_BLUE
               }
               filltext = 'S'
           }
           if (tchar == '2') {
               if (xchar == 'G') {
                   fill = MEDIUM_GREEN
               }
               if (xchar == 'B') {
                   fill = MEDIUM_BLUE
               }
               filltext = 'M'
           }
           if (tchar == '3') {
               if (xchar == 'G') {
                   fill = STOPLIGHT_GREEN
               }
               if (xchar == 'B') {
                   fill = STOPLIGHT_BLUE
               }
               filltext = 'L'
           }
           stroke = fill
       }
       paper.circle(ORIGINX+46*col, ORIGINY+46*row, 20).attr({ "stroke": stroke, "fill": fill });
       paper.text(ORIGINX+46*col, ORIGINY+46*row, filltext)
    }
   }
  },'html');
 }

 $.post('work', {}, function(r) {
  //$('#workresults').html('This is urlBar.gsps work plus work result: ' + r);
 },'html');

 var paper = Raphael(0, 0, 1200, 1200);
 console.log('xxx');
 si1 = setInterval(getMessages2, 200);

  </script>
  <!--
  <div id="workresults" >This is #workresults in urlBar.gsp...</div>
  -->
  <div id="controllerresults" ></div>
 </body>`

	w.Write([]byte(fmt.Sprintf(html)))

}

func work(w http.ResponseWriter, r *http.Request) {
	var piazzaBox string
	if containers != nil {
		piazzaBox = containers[0]
	} else {
		piazzaBox = myIPWithTimeout()
	}

	workers := 0 /* for now, don't use this or multiple invocations (no go global var) */

	//work() is invoked from the gsp page each time the browser is refreshed.
	//we don't want multiple workers, so each worker gets a worker number.
	//if the worker number doesn't match the number of workers, then
	//the worker exits.
	workers++
	iamworker := workers
	for i := 0; i < numColorfulDisplayDots; i++ {
		q[i] = newTestVector(piazzaBox, randomExternalUserService())
	}

	var maxIterationToCallTestVector = 64000
	for j := 0; j < maxIterationToCallTestVector; j++ {
		if iamworker == workers && !booleanOfDotCompletion() {
			var healthCheckServicesEverySoOften = 10
			if j%healthCheckServicesEverySoOften == 0 {
				qhealth = updateHealthArray(*qhealth)
			}
			for i := 0; i < numColorfulDisplayDots; i++ {
				q[i] = nextStep(*q[i])
			}
		}
	}
}

func myIPWithTimeout() string {
	timeout := time.Duration(3 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get("http://169.254.169.254/latest/meta-data/public-ipv4")
	if err != nil {
		fmt.Println("Something went wrong in myIPWithTimeout")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body)
}

func status(w http.ResponseWriter, r *http.Request) {
	var s1 string
	var s2 string
	var s4 string
	var s5 string
	if q[0] == nil {
		s1 = `1124012411241124`
		s2 = `1122331122331122`
		s4 = `hello`
		s5 = `hello`
	} else {
		s1 = stringOfDotStatusEachRepresentsAPiazzaJob()
		s2 = stringOfDotDurationEachRepresentsAPiazzaJob()
		s4 = statusString(qhealth)
		s5 = stringOfDotColor()
	}
	s := `{"dotStatus":"` + s1 + `",` +
		`"dotDuration":"` + s2 + `",` +
		`"squareHealth":"` + s4 + `",` +
		`"dotColor":"` + s5 + `"}`

	byu := []byte(s)
	w.Write(byu)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/work", work)
	mux.HandleFunc("/status", status)
	//mux.HandleFunc("/external", simulatedExternalUserService)
	http.ListenAndServe(":8076", mux)
}

/*
   def dots() {
   }
   def zdots() {
       zwork()
       string()
   }
   def string() {
       render stringOfDotStatusEachRepresentsAPiazzaJob() + '\n'
   }
*/
/* s/m: reinstate these four lines into pz1curl,...pz4curl
   timeout := time.Duration(3 * time.Second)
   client := http.Client {
       Timeout: timeout,
   }
*/
