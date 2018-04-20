var boxplotMode = false;

function createHistogram(hist, name, maxY, rspLength) {
  var scale = 1.0;
  if (rspLength > 1) {
    var histRange = hist[hist.length-1].V - hist[0].V
    scale = hist.length / histRange;
  }
  hist.forEach(h => {
    if (h.P*scale > maxY) maxY = h.P*scale;
  });
  return {
    x: hist.map(function(h) { return h.V }),
    y: hist.map(function(h) { return h.P * scale }),
    mode: 'lines',
    name: name
  };
}

function createBoxPlot(hist, name) {
  var nbSamples = 10000;
  var pCounter = 0.0;
  var sampleCounter = 0;
  return {
    x: hist.reduce(function(acc, h) { 
      while ((pCounter + h.P) >+ (sampleCounter/nbSamples)) {
        acc.push(h.V);
        sampleCounter++;
      }
      pCounter += h.P;
      return acc;
    }, []),
    type: 'box',
    boxmean: true,
    name: name
  };
}

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
            acc.push(boxplotMode ? createBoxPlot(hist, codes[i]) : createHistogram(hist, codes[i], maxY, response.length));
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

function onClickBoxplot(e) {
  boxplotMode = e.target.checked;
  generateDiagram();
}

document.getElementById("code").addEventListener("keyup", generateDiagram);
document.getElementById("boxplot").addEventListener("click", onClickBoxplot);
generateDiagram();