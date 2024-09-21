/*window.addEventListener('scroll', function(){

  var nav = document.getElementById('MainNav');
  var yPosition = window.pageYOffset;

  if(yPosition > 10) {
    nav.style.display='none';
  }else{
    nav.style.display='block';
  }

});*/

//On scroll stop show nav?

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
/*
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
*/

