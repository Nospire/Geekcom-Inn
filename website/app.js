const lines = [
  "pull up a chair • warm lights • low music",
  "the tavern is open now • the website follows soon",
  "ssh tavrn.sh • no signup • no browser",
];

const node = document.getElementById("status-line");
let index = 0;

setInterval(() => {
  index = (index + 1) % lines.length;
  node.textContent = lines[index];
}, 2400);
