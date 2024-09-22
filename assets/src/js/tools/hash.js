var str = document.getElementById('InputString');

str.addEventListener('keyup', function(){
    
  var iterations = document.getElementById('InputCount').value;

  var hashes = {}

  var sha512 = new jsSHA("SHA-512", "TEXT");
  var sha256 = new jsSHA("SHA-256", "TEXT");
  var sha1 = new jsSHA("SHA-1", "TEXT");
  
  for(let i = 0; i<iterations; i++){
    if(hash512==undefined){
        sha512.update(str.value);
        var hash512 = sha512.getHash("HEX");
    }else{
        sha512.update(hash512);
        hash512 = sha512.getHash("HEX");
    }

    if(hash256==undefined){
        sha256.update(str.value);
        var hash256 = sha256.getHash("HEX");
    }else{
        sha256.update(hash256);
        hash256 = sha256.getHash("HEX");
    }

    if(hash1==undefined){
        sha1.update(str.value);
        var hash1 = sha1.getHash("HEX");
    }else{
        sha1.update(hash1);
        hash1 = sha1.getHash("HEX");
    }
    
    if(hashMD5==undefined){
        var hashMD5 = md5(str.value);
    }else{
        hashMD5 = md5(hashMD5);
    }

    hashes.sha512 = hash512;
    hashes.sha256 = hash256;
    hashes.sha1 = hash1;
    hashes.md5 = hashMD5;
    
  }

  var hashText;
  for(let h in hashes){

    if(hashText == undefined){
        hashText=`<div><strong>${h}</strong>
            <div style="word-wrap:break-word;">${hashes[h]}</div>
        </div>`
    }else{
        hashText = hashText + `<hr><div><strong>${h}</strong>
            <div style="word-wrap:break-word;">${hashes[h]}</div>
        </div>`
    }

  }

  document.getElementById('Hashes').innerHTML = hashText;
});