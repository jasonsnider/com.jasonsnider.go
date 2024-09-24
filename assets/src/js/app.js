document.getElementById('ShowMainNav').addEventListener('click',function(){

    var nav = document.querySelector('nav#MainNav ul');
  
    if (nav.style.display === "block") {
        nav.style.display = "none";
    } else {
        nav.style.display = "block";
    }
  
  });