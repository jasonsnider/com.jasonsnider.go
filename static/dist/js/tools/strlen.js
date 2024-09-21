var str = document.getElementById('InputString');

str.addEventListener('keyup', function(){
  document.getElementById('StringLength').innerHTML = str.value.length;
  document.getElementById('Results').style.display = 'inline';
});
