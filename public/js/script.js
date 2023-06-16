let blogs = [];

function getProject(event) {
  event.preventDefault();

  //get id from html
  let projectName = document.getElementById("projectName").value;
  let startDate = document.getElementById("startDate").value;
  let endDate = document.getElementById("endDate").value;
  let desc = document.getElementById("desc").value;
  let image = document.getElementById("upload-image").files;

  // declare icon technologies
  const tech1 = '<i class="fa-brands fa-node"></i>';
  const tech2 = '<i class="fa-brands fa-react"></i>';
  const tech3 = '<i class="fa-brands fa-python"></i>';
  const tech4 = '<i class="fa-brands fa-golang"></i>';

  //check icon
  let nodeJS = document.getElementById("tech-one").checked ? tech1 : "";
  let reactJS = document.getElementById("tech-two").checked ? tech2 : "";
  let python = document.getElementById("tech-three").checked ? tech3 : "";
  let golang = document.getElementById("tech-four").checked ? tech4 : "";

  //validation
  if (projectName == "") {
    return alert("Tolong isi bagian nama projek anda");
  } else if (startDate == "") {
    return alert("Tolong masukkan tanggal yang sesuai");
  } else if (endDate == "") {
    return alert("Tolong masukkan tanggal yang sesuai");
  } else if (desc == "") {
    return alert("Tolong masukkan deskripsi dengan benar");
  } else {
    alert("Segera ditampilkan");
  }

  // convert image to blob object
  image = URL.createObjectURL(image[0]);
  console.log(image);

  //count duration
  startDate = new Date(startDate);
  endDate = new Date(endDate);

  const distance = endDate.getTime() - startDate.getTime();

  let day = distance / (1000*3600*24);
  if(day > 0) day=`${day} hari`


  let dataProject = {
    image,
    projectName,
    startDate,
    endDate,
    day,
    desc,
    nodeJS,
    reactJS,
    python,
    golang,
    author: "Ilham Setyadji",
    postAt: new Date(),
  };

  blogs.push(dataProject);
  console.log(blogs);

  clear();
  showProject();
}

function showProject() {

  //menghapus data post
  document.getElementById("conten").innerHTML = "";

  //menambah data post
  for (let i = 0; i < blogs.length; i++) {
    document.getElementById("conten").innerHTML += `
              <div class="card-project">
                <div class="content">
                  <img src="${blogs[i].image}" alt="image">
                  <a href="add-project-detail.html" target="_blank"><h4>${blogs[i].projectName}</h4></a>
                  <div class="detail-blog-content">durasi ${blogs[i].day} | ${blogs[i].author} </div>
                    <div class="text">
                      <p>${blogs[i].desc}</p>
                    </div>
                    <!-- technologies -->
                    <div class="tech">
                    ${blogs[i].nodeJS}
                    ${blogs[i].reactJS}
                    ${blogs[i].python}
                    ${blogs[i].golang}
                    </div>
                    <!-- end-technologies -->
                            
                </div>
                    
                <button>edit</button>
                <button>delete</button>
                <div>
                <p style="font-size: 15px; float: right; color: grey">${getDistanceTime(blogs[i].postAt)}</p>
                </div>
              </div>
              `;
  }
}

//menghapus form input setelah diisi
function clear() {
  document.getElementById("projectName").value = "";
  document.getElementById("startDate").value = "";
  document.getElementById("endDate").value = "";
  document.getElementById("desc").value = "";
  document.getElementById("upload-image").value = "";
}



function getFullTime(time) {
  // console.log("get full time");
  // let time = new Date();
  // console.log(time);

  let monthName = [
    "Jan",
    "Feb",
    "Mar",
    "Apr",
    "May",
    "Jun",
    "Jul",
    "Aug",
    "Sep",
    "Oct",
    "Nov",
    "Dec",
  ];
  // console.log(monthName[8]);

  let date = time.getDate();
  // console.log(date);

  let monthIndex = time.getMonth();
  // console.log(monthIndex);

  let year = time.getFullYear();
  // console.log(year);

  let hours = time.getHours();
  let minutes = time.getMinutes();
  // console.log(minutes);

  if (hours <= 9) {
    hours = "0" + hours;
  } else if (minutes <= 9) {
    minutes = "0" + minutes;
  }

  return `${date} ${monthName[monthIndex]} ${year} ${hours}:${minutes} WIB`;
}

function getDistanceTime(time) {
  let timeNow = new Date();
  let timePost = time;

  // waktu sekarang - waktu post
  let distance = timeNow - timePost; // hasilnya milidetik
  console.log(distance);

  let milisecond = 1000; // milisecond
  let secondInHours = 3600; // 1 jam 3600 detik
  let hoursInDays = 24; // 1 hari 24 jam

  let distanceDay = Math.floor(distance / (milisecond * secondInHours * hoursInDays)); // 1/86400000
  let distanceHours = Math.floor(distance / (milisecond * 60 * 60)); // 1/3600000
  let distanceMinutes = Math.floor(distance / (milisecond * 60)); // 1/60000
  let distanceSeconds = Math.floor(distance / milisecond); // 1/1000

  if (distanceDay > 0) {
    return `${distanceDay} Day Ago`;
  } else if (distanceHours > 0) {
    return `${distanceHours} Hours Ago`;
  } else if (distanceMinutes > 0) {
    return `${distanceMinutes} Minutes Ago`;
  } else {
    return `${distanceSeconds} Seconds Ago`;
  }
}

setInterval(function () {
  showProject();
}, 50000);

