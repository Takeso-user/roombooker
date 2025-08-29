// Room Booker Frontend

let calendar;
let selectedRoomId = null;
let currentUser = null;

document.addEventListener("DOMContentLoaded", function () {
  initAuth();
  initCalendar();
  loadOffices();
});

function initAuth() {
  const loginBtn = document.getElementById("login-btn");
  loginBtn.addEventListener("click", () => {
    window.location.href = "/auth/oidc/start";
  });

  // Check if user is logged in
  fetch("/me", {
    headers: {
      Authorization: `Bearer ${localStorage.getItem("token")}`,
    },
  })
    .then((response) => {
      if (response.ok) {
        return response.json();
      }
      throw new Error("Not authenticated");
    })
    .then((user) => {
      currentUser = user;
      document.getElementById("auth-section").style.display = "none";
      document.getElementById("main-content").style.display = "block";
    })
    .catch(() => {
      // Not logged in
    });
}

function initCalendar() {
  const calendarEl = document.getElementById("calendar");
  calendar = new FullCalendar.Calendar(calendarEl, {
    initialView: "timeGridWeek",
    headerToolbar: {
      left: "prev,next today",
      center: "title",
      right: "dayGridMonth,timeGridWeek,timeGridDay",
    },
    selectable: true,
    select: function (info) {
      if (selectedRoomId) {
        openBookingModal(info.start, info.end);
      }
    },
    eventClick: function (info) {
      // Handle event click
    },
    events: function (fetchInfo, successCallback, failureCallback) {
      if (!selectedRoomId) {
        successCallback([]);
        return;
      }
      fetch(
        `/rooms/${selectedRoomId}/calendar?from=${fetchInfo.start.toISOString()}&to=${fetchInfo.end.toISOString()}`,
        {
          headers: {
            Authorization: `Bearer ${localStorage.getItem("token")}`,
          },
        }
      )
        .then((response) => response.json())
        .then((events) => successCallback(events))
        .catch((error) => failureCallback(error));
    },
  });
  calendar.render();
}

function loadOffices() {
  fetch("/offices", {
    headers: {
      Authorization: `Bearer ${localStorage.getItem("token")}`,
    },
  })
    .then((response) => response.json())
    .then((offices) => {
      const select = document.getElementById("office-select");
      offices.forEach((office) => {
        const option = document.createElement("option");
        option.value = office.id;
        option.textContent = office.name;
        select.appendChild(option);
      });
    });
}

function loadFloors(officeId) {
  fetch(`/floors?office_id=${officeId}`, {
    headers: {
      Authorization: `Bearer ${localStorage.getItem("token")}`,
    },
  })
    .then((response) => response.json())
    .then((floors) => {
      const select = document.getElementById("floor-select");
      select.innerHTML = '<option value="">Select Floor</option>';
      floors.forEach((floor) => {
        const option = document.createElement("option");
        option.value = floor.id;
        option.textContent = floor.label;
        select.appendChild(option);
      });
    });
}

function loadRooms(floorId) {
  fetch(`/rooms?floor=${floorId}`, {
    headers: {
      Authorization: `Bearer ${localStorage.getItem("token")}`,
    },
  })
    .then((response) => response.json())
    .then((rooms) => {
      const list = document.getElementById("room-list");
      list.innerHTML = "";
      rooms.forEach((room) => {
        const div = document.createElement("div");
        div.className = "form-check";
        div.innerHTML = `
                <input class="form-check-input" type="radio" name="room" id="room-${room.id}" value="${room.id}">
                <label class="form-check-label" for="room-${room.id}">
                    ${room.name} (Capacity: ${room.capacity})
                </label>
            `;
        div.querySelector("input").addEventListener("change", (e) => {
          if (e.target.checked) {
            selectedRoomId = e.target.value;
            calendar.refetchEvents();
          }
        });
        list.appendChild(div);
      });
    });
}

function openBookingModal(start, end) {
  document.getElementById("start").value = start.toISOString().slice(0, 16);
  document.getElementById("end").value = end.toISOString().slice(0, 16);
  const modal = new bootstrap.Modal(document.getElementById("bookingModal"));
  modal.show();
}

document.getElementById("office-select").addEventListener("change", (e) => {
  if (e.target.value) {
    loadFloors(e.target.value);
  }
});

document.getElementById("floor-select").addEventListener("change", (e) => {
  if (e.target.value) {
    loadRooms(e.target.value);
  }
});

document.getElementById("save-booking").addEventListener("click", () => {
  const form = document.getElementById("booking-form");
  const formData = new FormData(form);
  const data = {
    title: formData.get("title"),
    start: new Date(formData.get("start")).toISOString(),
    end: new Date(formData.get("end")).toISOString(),
    participants: formData
      .get("participants")
      .split("\n")
      .filter((email) => email.trim()),
  };

  fetch(`/rooms/${selectedRoomId}/bookings`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${localStorage.getItem("token")}`,
    },
    body: JSON.stringify(data),
  }).then((response) => {
    if (response.ok) {
      calendar.refetchEvents();
      bootstrap.Modal.getInstance(
        document.getElementById("bookingModal")
      ).hide();
      form.reset();
    } else {
      alert("Failed to create booking");
    }
  });
});
