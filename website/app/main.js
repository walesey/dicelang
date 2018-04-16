function generateDiagram(){
  var codes = document.getElementById("code").value.split("\n")
  var code = JSON.stringify(codes);
  var xhttp = new XMLHttpRequest();
  xhttp.onreadystatechange = function() {
    if (this.readyState == 4) {
      if (this.status == 200) {
        var response = JSON.parse(this.responseText);

        var layout = {
          xaxis: {
            title: 'Result',
          },
          yaxis: {
            title: 'Probability',
          }
        };
        
        var maxY = 0;
        var data = response.reduce(function(acc, hist, i) {
          if (typeof hist != 'number'){
            hist.forEach(h => {
              if (h.P > maxY) maxY = h.P;
            });
            acc.push({
              x: hist.map(function(h) { return h.V }),
              y: hist.map(function(h) { return h.P }),
              mode: 'lines',
              name: codes[i]
            });
          }
          return acc;
        }, []);

        response.forEach(function(hist, i) {
          if (typeof hist == 'number'){
            data.push({
              x: [hist, hist],
              y: [0.0, maxY],
              mode: 'lines',
              name: codes[i]
            })
          }
        });

        console.log(data);

        document.getElementById("plot").innerHTML = "";
        Plotly.newPlot('plot', data, layout);
      } else {
        document.getElementById("plot").innerHTML = this.responseText;
      }
    }
  };
  xhttp.open("POST", "/code", true);
  xhttp.send(code);

  var b64Code = btoa(document.getElementById("code").value);
  document.getElementById("share").setAttribute("href", "?code=" + b64Code);
}

document.getElementById("code").addEventListener("keyup", generateDiagram);
generateDiagram();