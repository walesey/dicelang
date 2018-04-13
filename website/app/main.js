function generateDiagram(){
  var code = document.getElementById("code").value;
  var xhttp = new XMLHttpRequest();
  xhttp.onreadystatechange = function() {
    if (this.readyState == 4) {
      if (this.status == 200) {
        var data = JSON.parse(this.responseText);

        var layout = {
          title: code,
          xaxis: {
            title: 'Result',
          },
          yaxis: {
            title: 'Probability',
          }
        };
        
        var data = {
          x: data.map(function(d) { return d.V }),
          y: data.map(function(d) { return d.P }),
          mode: 'lines',
          name: 'Lines'
        };
        
        document.getElementById("plot").innerHTML = "";
        Plotly.newPlot('plot', [data], layout);
      } else {
        document.getElementById("plot").innerHTML = this.responseText;
      }
    }
  };
  xhttp.open("POST", "/code", true);
  xhttp.send(code);

  document.getElementById("share").setAttribute("href", "?code=" + btoa(code))
}

document.getElementById("code").addEventListener("keyup", generateDiagram);
generateDiagram();