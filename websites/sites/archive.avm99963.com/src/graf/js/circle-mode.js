circleMode = false;

function initCircleMode() {
  document.querySelector("#circle-mode").addEventListener('click', function() {
    circleMode = !circleMode;
    document.querySelector("#circle-mode i").innerText = (circleMode ? "scatter_plot" : "trip_origin");
    s.graph.nodes().forEach(function(n) {
      n.x = (circleMode ? n.circleX : n.originalX);
      n.y = (circleMode ? n.circleY : n.originalY);
      n.size = 10;
    });

    s.refresh();
  });
}
