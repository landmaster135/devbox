// script.js - Image Loading and Mouse Tracking Bar Functionality

const handleImageChange = (e) => {
  update(false, { isInputting: false, e: null }, { isUpdating: false, e: e }, { isImageChanged: true, e: e });
  updateLockStatusDisplay(true, false); // Show default lock state post-image loading
}
document.getElementById("imageLoader").addEventListener("change", handleImageChange);

const handleBarsAndCoordinates = (e) => {
  const isLocked = updateLockStatusDisplay(false, null); // Update lock status on mouse move
  if (!isLocked) { // Allow bar movement only if not locked
    update(false, { isInputting: false, e: null }, { isUpdating: true, e: e }, { isImageChanged: false, e: null });
  }
}
document.addEventListener("mousemove", handleBarsAndCoordinates);

const updateCoordinateInput = (e) => {
  updateLockStatusDisplay(true, true); // Update lock status
  update(false, { isInputting: true, e: e }, { isUpdating: false, e: e }, { isImageChanged: false, e: null });
}
document.addEventListener("mousedown", updateCoordinateInput);

const unlockBar = (e) => {
  updateLockStatusDisplay(true, false); // Update lock status
  update(true, { isInputting: false, e: null }, { isUpdating: false, e: null }, { isImageChanged: false, e: null });
}
document.addEventListener("mouseup", unlockBar);

const updateLockStatusDisplay = (isUpdating, isLocking) => {
  let flg;
  if (isUpdating) {
    flg = isLocking;
  }

  const lockStatusDisplay = document.getElementById("lockStatus"); // element to indicate lock state
  lockStatusDisplay.innerText = flg ? "Locking..." : "Releasing now.";
  lockStatusDisplay.style.color = flg ? "red" : "green"; // Change color based on lock status
  return flg;
}

const copyToClipboard = (copyStatusesToCb, coordinateXUpper, coordinateYUpper, coordinateXLower, coordinateYLower) => {
  const sep = " ";
  // const copyText = `${coordinateXUpper.value}${sep}${coordinateYUpper.value}${sep}${coordinateXLower.value}${sep}${coordinateYLower.value}`;
  const copyText = `-x1 ${coordinateXUpper.value}${sep}-y1 ${coordinateYUpper.value}${sep}-x2 ${coordinateXLower.value}${sep}-y2 ${coordinateYLower.value}`;

  // クリップボード表示用テキストボックスに座標を表示
  const clipboardCoordinates = document.getElementById("clipboardCoordinates");
  clipboardCoordinates.value = copyText;

  // クリップボードにコピー
  navigator.clipboard.writeText(copyText).then(() => {
    const newElement = document.createElement("p");
    newElement.textContent = "copied!";
    copyStatusesToCb.insertAdjacentElement("afterend", newElement);
    setTimeout(() => {
      newElement.remove();
    }, 1000);
  }, (err) => {
    console.error("Asynchronous clipboard write failed.", err);
  });
}

const updateStateOfBarActivation = (indexOfActiveBar) => {
  const NUMBER_OF_BAR = 2;
  return (isBarStateUpdating, propsForCoordRecordInput, propsForUpdateCoordAndBar, propsForImageChanging) => {
    console.log("updateStateOfBarActivation(): updating state start.");
    if (isBarStateUpdating) {
      indexOfActiveBar += 1;
      if (indexOfActiveBar === NUMBER_OF_BAR) {
        indexOfActiveBar = 0;
      }
    }
    console.log("updateStateOfBarActivation(): the active bar is decided.");

    const imageElement = document.getElementById("imageElement");
    let coordinateXUpperToRecord = document.getElementById("coordinateXUpperDisplayToRecord");
    let coordinateYUpperToRecord = document.getElementById("coordinateYUpperDisplayToRecord");
    let coordinateXLowerToRecord = document.getElementById("coordinateXLowerDisplayToRecord");
    let coordinateYLowerToRecord = document.getElementById("coordinateYLowerDisplayToRecord");
    const copyStatusDisplayToCb = document.getElementById("copyToCb");

    if (propsForCoordRecordInput.isInputting) {
      const imageRect = imageElement.getBoundingClientRect();
      const x = Math.round(propsForCoordRecordInput.e.clientX - imageRect.left);
      const y = Math.round(propsForCoordRecordInput.e.clientY - imageRect.top);
      switch (indexOfActiveBar) {
        case 0:
          coordinateXUpperToRecord.value = `${x}`;
          coordinateYUpperToRecord.value = `${y}`;
          break;
        case 1:
          coordinateXLowerToRecord.value = `${x}`;
          coordinateYLowerToRecord.value = `${y}`;
          break;
      }
      copyToClipboard(copyStatusDisplayToCb, coordinateXUpperToRecord, coordinateYUpperToRecord, coordinateXLowerToRecord, coordinateYLowerToRecord);
    }
    console.log("updateStateOfBarActivation(): coordinates are recorded.");

    const coordinateXUpperDisplay = document.getElementById("coordinateXUpperDisplay");
    const coordinateYUpperDisplay = document.getElementById("coordinateYUpperDisplay");
    const coordinateXLowerDisplay = document.getElementById("coordinateXLowerDisplay");
    const coordinateYLowerDisplay = document.getElementById("coordinateYLowerDisplay");
    const draggableBarXUpper = document.getElementById("draggableBarXUpper");
    const draggableBarYUpper = document.getElementById("draggableBarYUpper");
    const draggableBarXLower = document.getElementById("draggableBarXLower");
    const draggableBarYLower = document.getElementById("draggableBarYLower");

    if (propsForUpdateCoordAndBar.isUpdating) {
      const updateBarsAndCoordinates = (imgRect, ev, barX, barY, coordX, coordY) => {
        if (ev.clientX >= imgRect.left && ev.clientX <= imgRect.right &&
          ev.clientY >= imgRect.top && ev.clientY <= imgRect.bottom) {
          // Track the bar to the mouse within the image
          const yPos = ev.clientY - imgRect.top;
          barY.style.top = `${yPos}px`;
          barY.style.transform = 'translateY(-50%)'; // Always center align
          coordY.value = `${Math.round(yPos)}`;

          const xPos = ev.clientX - imgRect.left;
          barX.style.left = `${xPos}px`;
          coordX.value = `${Math.round(xPos)}`;
        }
      }
      const imageRect = imageElement.getBoundingClientRect();
      let draggableBarX = draggableBarXUpper;
      let draggableBarY = draggableBarYUpper;
      let coordXDisplay = coordinateXUpperDisplay;
      let coordYDisplay = coordinateYUpperDisplay;
      switch (indexOfActiveBar) {
        case 0:
          // Nothing to do.
          break;
        case 1:
          draggableBarX = draggableBarXLower;
          draggableBarY = draggableBarYLower;
          coordXDisplay = coordinateXLowerDisplay;
          coordYDisplay = coordinateYLowerDisplay;
          break;
      }
      updateBarsAndCoordinates(imageRect, propsForUpdateCoordAndBar.e, draggableBarX, draggableBarY, coordXDisplay, coordYDisplay);
    }
    console.log("updateStateOfBarActivation(): bars and coordinates are updated.");

    console.log(propsForImageChanging)
    if (propsForImageChanging.isImageChanged) {
      const selectedImage = propsForImageChanging.e.target.files[0];
      if (!selectedImage) return;
      const reader = new FileReader();
      reader.onload = (event) => {
        imageElement.src = event.target.result;
        // Initialize bar after image load
        draggableBarYUpper.style.top = "50%"; // Align to image center
        draggableBarYUpper.style.transform = "translateY(-50%)";
        draggableBarXUpper.style.left = "0px"; // Align to image left edge
        coordinateXUpperDisplay.value = "0";
        coordinateYUpperDisplay.value = `${Math.round(imageElement.offsetHeight / 2)}`;
      };
      reader.readAsDataURL(selectedImage);
    }
    console.log("updateStateOfBarActivation(): image is loaded.");

    console.log("updateStateOfBarActivation(): updating state finished.")
  }
}

const update = updateStateOfBarActivation(0);
