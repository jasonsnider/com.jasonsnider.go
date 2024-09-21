if(document.getElementById('MainNav')!==null){

  window.addEventListener(
    'scroll',

    function(){

      var nav = document.getElementById('MainNav');
      var yPosition = window.pageYOffset;

      if(yPosition > 10) {
        nav.style.display='block';
      }else{
        nav.style.display='none';
      }

    }

  );

  document.getElementById('ToAbout').addEventListener(
    'click',
    function(e) {


      var el = document.getElementById('About');
      var rect = el.getBoundingClientRect();
      window.scrollTo({
          top: (rect.top + window.pageYOffset) - 60,
          behavior: "smooth"
      });

      history.pushState(null, null, '#ToAbout');
      e.preventDefault();
    }
  );

  document.getElementById('ToAboutNext').addEventListener(
    'click',
    function(e) {


      var el = document.getElementById('About');
      var rect = el.getBoundingClientRect();

      window.scrollTo({
          top: (rect.top + window.pageYOffset) - 60,
          behavior: "smooth"
      });

      history.pushState(null, null, '#ToAbout');
      e.preventDefault();
    }
  );

  document.getElementById('ToProject').addEventListener(
    'click',
    function(e) {


      var el = document.getElementById('Project');
      var rect = el.getBoundingClientRect();
      window.scrollTo({
          top: (rect.top + window.pageYOffset) - 60,
          behavior: "smooth"
      });

      history.pushState(null, null, '#ToProject');
      e.preventDefault();
    }
  );

}


/* Highlight.js */
window.onload = function() {
  var aCodes = document.getElementsByTagName('pre');
  for (var i=0; i < aCodes.length; i++) {
      hljs.highlightBlock(aCodes[i]);
  }
};
/* Load Defered styles */
/*
var loadDeferredStyles = function() {
   var addStylesNode = document.getElementById("deferredStyles");
   var replacement = document.createElement("div");
   replacement.innerHTML = addStylesNode.textContent;
   document.body.appendChild(replacement)
   addStylesNode.parentElement.removeChild(addStylesNode);
};

var raf = requestAnimationFrame || mozRequestAnimationFrame ||
  webkitRequestAnimationFrame || msRequestAnimationFrame;

if (raf) raf(function() { window.setTimeout(loadDeferredStyles, 0); });
else window.addEventListener('load', loadDeferredStyles);
*/
