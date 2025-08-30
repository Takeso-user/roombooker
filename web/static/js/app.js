// Room Booker Frontend with Microsoft-style UI

let calendar;
let selectedRoomId = null;
let currentUser = null;

document.addEventListener("DOMContentLoaded", function () {
  initAuth();
  initCalendar();
  loadOffices();
});

function initAuth() {
  // Check if user is logged in
  checkAuthStatus();
}

function checkAuthStatus() {
  fetch("/me", {
    credentials: "same-origin",
  })
    .then((response) => {
      if (response.ok) {
        return response.json();
      }
      throw new Error("Not authenticated");
    })
    .then((user) => {
      currentUser = user;
      updateUserHeader(user);
      loadOffices();
      if (user.role === "admin") {
        showAdminPanel();
      }
    })
    .catch(() => {
      // Not logged in, redirect to login
      window.location.href = "/login";
    });
}

function updateUserHeader(user) {
  document.getElementById("userName").textContent = user.name;
  document.getElementById("userRole").textContent = user.role;
  document.getElementById("userAvatar").textContent = user.name
    .charAt(0)
    .toUpperCase();

  // Add admin button if user is admin
  if (user.role === "admin") {
    const adminBtn = document.createElement("button");
    adminBtn.className = "btn btn-outline-primary btn-sm ms-2";
    adminBtn.textContent = "Admin Panel";
    adminBtn.onclick = toggleAdminPanel;
    document.getElementById("userInfo").appendChild(adminBtn);
  }
}

function logout() {
  fetch("/auth/logout", {
    method: "POST",
    credentials: "same-origin",
  }).then(() => {
    window.location.href = "/login";
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
      // Handle event click - could show booking details
      console.log("Event clicked:", info.event);
    },
    events: function (fetchInfo, successCallback, failureCallback) {
      if (!selectedRoomId) {
        successCallback([]);
        return;
      }
      fetch(
        `/api/rooms/${selectedRoomId}/bookings?from=${fetchInfo.start.toISOString()}&to=${fetchInfo.end.toISOString()}`,
        {
          credentials: "same-origin",
        }
      )
        .then((response) => response.json())
        .then((events) => {
          console.debug(
            "Fetched events for room",
            selectedRoomId,
            "count",
            events.length,
            events
          );
          successCallback(events);
        })
        .catch((error) => failureCallback(error));
    },
  });
  calendar.render();
}

function loadOffices() {
  fetch("/api/offices", {
    credentials: "same-origin",
  })
    .then((response) => response.json())
    .then((offices) => {
      const select = document.getElementById("officeSelect");
      select.innerHTML = '<option value="">Select Office</option>';
      offices.forEach((office) => {
        const option = document.createElement("option");
        option.value = office.id;
        option.textContent = office.name;
        select.appendChild(option);
      });
    });
}

function loadRooms() {
  const officeId = document.getElementById("officeSelect").value;
  if (!officeId) return;

  fetch(`/api/offices/${officeId}/rooms`, {
    credentials: "same-origin",
  })
    .then((response) => response.json())
    .then((rooms) => {
      const select = document.getElementById("roomSelect");
      select.innerHTML = '<option value="">Select Room</option>';
      rooms.forEach((room) => {
        const option = document.createElement("option");
        option.value = room.id;
        option.textContent = `${room.name} (Capacity: ${room.capacity})`;
        select.appendChild(option);
      });
      console.debug("Rooms loaded for office", officeId, rooms.length);
    });
}

function loadCalendar() {
  selectedRoomId = document.getElementById("roomSelect").value;
  if (selectedRoomId) {
    calendar.refetchEvents();
  }
}

function openBookingModal(start, end) {
  document.getElementById("startTime").value = start.toISOString().slice(0, 16);
  document.getElementById("endTime").value = end.toISOString().slice(0, 16);
  const modal = new bootstrap.Modal(document.getElementById("bookingModal"));
  modal.show();
}

function saveBooking() {
  const form = document.getElementById("bookingForm");

  // Validate client-side: room must be selected
  if (!selectedRoomId) {
    showErrorMessage("Please select a room before booking.");
    return;
  }

  const data = {
    title: document.getElementById("title").value,
    start_time: document.getElementById("startTime").value,
    end_time: document.getElementById("endTime").value,
    attendees: document
      .getElementById("attendees")
      .value.split("\n")
      .filter((email) => email.trim()),
    room_id: selectedRoomId,
  };

  fetch(`/api/bookings`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "same-origin",
    body: JSON.stringify(data),
  })
    .then(async (response) => {
      if (response.ok) {
        // parse created booking and log for debugging
        let created = null;
        try {
          created = await response.json();
          console.debug("Booking created:", created);
        } catch (e) {
          console.debug("Booking created but response not JSON");
        }

        calendar.refetchEvents();
        const modal = bootstrap.Modal.getInstance(
          document.getElementById("bookingModal")
        );
        if (modal) modal.hide();
        form.reset();
        // If we have created booking payload, add it to calendar immediately and navigate to it
        if (created && created.start) {
          try {
            calendar.addEvent({
              id: created.id,
              title: created.title,
              start: created.start,
              end: created.end,
              color: created.color || undefined,
            });
            // ensure visible
            // make sure the room is selected and calendar shows the event
            selectedRoomId = created.room_id || selectedRoomId;
            const roomSel = document.getElementById("roomSelect");
            if (roomSel) {
              try {
                roomSel.value = selectedRoomId;
              } catch (e) {}
            }
            calendar.refetchEvents();
            calendar.gotoDate(created.start);
          } catch (e) {
            console.debug("Failed to add event to calendar immediately", e);
          }
        }
        showSuccessMessage(
          "Booking created successfully! " +
            (created ? `ID: ${created.id}` : "")
        );
        return;
      }

      // Try to parse JSON error, fallback to plain text
      let errText = "Failed to create booking";
      try {
        const j = await response.json();
        errText = j.message || JSON.stringify(j) || errText;
      } catch (e) {
        try {
          errText = await response.text();
        } catch (e2) {
          /* ignore */
        }
      }
      showErrorMessage(errText || `Booking failed (status ${response.status})`);
    })
    .catch((err) => {
      console.error("Booking error:", err);
      showErrorMessage("Network or server error while creating booking");
    });
}

// Admin Panel Functions
function showAdminPanel() {
  document.getElementById("adminPanel").style.display = "block";
}

function toggleAdminPanel() {
  const panel = document.getElementById("adminPanel");
  panel.style.display = panel.style.display === "none" ? "block" : "none";
}

function showCreateRoomModal() {
  // Implementation for creating rooms
  alert("Create room functionality will be implemented");
}

function loadRoomsList() {
  // Implementation for loading all rooms
  alert("Load rooms list functionality will be implemented");
}

function showUserManagement() {
  // Implementation for user role management
  alert("User management functionality will be implemented");
}

function loadUsersList() {
  fetch("/api/admin/users", { credentials: "same-origin" })
    .then((r) => r.json())
    .then((users) => {
      const el = document.getElementById("usersList");
      el.innerHTML = "";
      users.forEach((u) => {
        const div = document.createElement("div");
        div.className = "user-row";
        div.textContent = `${u.email} — ${u.role}`;
        el.appendChild(div);
      });
    })
    .catch((err) => {
      console.error(err);
      showErrorMessage("Failed to load users");
    });
}

function showAdminSection(name) {
  document
    .querySelectorAll(".admin-section")
    .forEach((s) => (s.style.display = "none"));
  const el = document.getElementById("admin-" + name);
  if (el) el.style.display = "block";
}

function createOffice() {
  const name = document.getElementById("officeName").value;
  const tz = document.getElementById("officeTZ").value || "UTC";
  fetch("/api/admin/offices", {
    method: "POST",
    credentials: "same-origin",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ name, timezone: tz }),
  })
    .then((r) => r.json())
    .then((j) => {
      showSuccessMessage("Office created: " + j.name);
    })
    .catch((e) => {
      console.error(e);
      showErrorMessage("Failed to create office");
    });
}

function createFloor() {
  const officeId = document.getElementById("floorOfficeId").value;
  const num = parseInt(document.getElementById("floorNumber").value || "0", 10);
  const label = document.getElementById("floorLabel").value;
  fetch("/api/admin/floors", {
    method: "POST",
    credentials: "same-origin",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ office_id: officeId, number: num, label }),
  })
    .then((r) => r.json())
    .then((j) => {
      showSuccessMessage("Floor created");
    })
    .catch((e) => {
      console.error(e);
      showErrorMessage("Failed to create floor");
    });
}

function createRoom() {
  const floorId = document.getElementById("roomFloorId").value;
  const name = document.getElementById("roomName").value;
  const capacity = parseInt(
    document.getElementById("roomCapacity").value || "0",
    10
  );
  const equipment = document.getElementById("roomEquipment").value;
  fetch("/api/admin/rooms", {
    method: "POST",
    credentials: "same-origin",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ floor_id: floorId, name, capacity, equipment }),
  })
    .then((r) => r.json())
    .then((j) => {
      showSuccessMessage("Room created: " + j.name);
    })
    .catch((e) => {
      console.error(e);
      showErrorMessage("Failed to create room");
    });
}

function showSystemInfo() {
  // Implementation for system information
  alert("System info functionality will be implemented");
}

function clearCache() {
  // Implementation for clearing cache
  alert("Clear cache functionality will be implemented");
}

// Utility Functions
function showSuccessMessage(message) {
  // Simple alert for now, could be enhanced with toast notifications
  alert("✅ " + message);
}

function showErrorMessage(message) {
  alert("❌ " + message);
}
