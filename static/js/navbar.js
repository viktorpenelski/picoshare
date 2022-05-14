import { logOut } from "./controllers/auth.js";
import { getStorageSpaceOrNull } from "./controllers/storageSpace.js";

document.addEventListener("DOMContentLoaded", () => {
  const navbarBurgers = Array.prototype.slice.call(
    document.querySelectorAll(".navbar-burger"),
    0
  );

  if (navbarBurgers.length > 0) {
    navbarBurgers.forEach((el) => {
      el.addEventListener("click", () => {
        const target = document.getElementById(el.dataset.target);
        el.classList.toggle("is-active");
        target.classList.toggle("is-active");
      });
    });
  }
});

const logOutEl = document.getElementById("navbar-log-out");
if (logOutEl) {
  logOutEl.addEventListener("click", () => {
    logOut().then(() => {
      document.location = "/";
    });
  });
}

getStorageSpaceOrNull().then((storageSpaceData) => {
  const storageSpaceEl = document.getElementById("storage-space");
  if (!storageSpaceData) {
    storageSpaceEl.style.visibility = "hidden";
    return;
  }

  storageSpaceEl.style.visibility = "visible";
  const storageSpaceInfoEl = document.getElementById("storage-space-info");
  storageSpaceInfoEl.innerHTML = `Storage (<b>${storageSpaceData.usedPercentage}%</b> full)<br>
<b>${storageSpaceData.used}</b> GB of <b>${storageSpaceData.all}</b> GB used`
  return;
});
