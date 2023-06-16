function getData() {
    // get value from id in html
    let name = document.getElementById("name").value;
    let email = document.getElementById("email").value;
    let phoneNumber = document.getElementById("phone-number").value;
    let subject = document.getElementById("subject").value;
    let message = document.getElementById("msg").value;
  
    if (name == ""){
      return alert("isi nama anda")
    } else if (email == ""){
      return alert("isi email anda")
    } else if (phoneNumber == ""){
      return alert("isi no hp anda")
    } else if (subject == ""){
      return alert("pilih subject anda")
    } else if (message == ""){
      return alert ("isi pesan anda!")
    } 
  
    const emailReceive = "ilhamsetyadji28@gmail.com";
  
    let gmailTo = document.createElement("a");
    gmailTo.href = `mailto:${emailReceive}?subject=${subject}&body=Halo nama saya ${name}, saya ingin ${message}, bisakah anda menghubungi saya di ${phoneNumber}`;
    gmailTo.click();
  
    let flex = {
      name,
      email,
      phoneNumber,
      subject,
      message,
    };
    console.log(flex);
  }