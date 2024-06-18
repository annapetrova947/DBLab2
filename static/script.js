document.addEventListener("DOMContentLoaded", function() {
    const modal = document.getElementById("editModal");
    const span = document.getElementsByClassName("close")[0];
    const editForm = document.getElementById("editForm");

    // Обработка кнопок редактирования
    const editButtons = document.getElementsByClassName("edit-btn");
    Array.from(editButtons).forEach(button => {
        button.addEventListener("click", function() {
            const id = this.getAttribute("data-id");
            const row = this.parentElement.parentElement;
            document.getElementById("editId").value = id;
            document.getElementById("editCoordinates").value = row.cells[1].innerText;
            document.getElementById("editLightIntensity").value = row.cells[2].innerText;
            document.getElementById("editForeignObjects").value = row.cells[3].innerText;
            document.getElementById("editStarObjectsCount").value = row.cells[4].innerText;
            document.getElementById("editUnknownObjectsCount").value = row.cells[5].innerText;
            document.getElementById("editDefinedObjectsCount").value = row.cells[6].innerText;
            document.getElementById("editNotes").value = row.cells[7].innerText;
            document.getElementById("editObjectID").value = row.cells[8].innerText;
            document.getElementById("editType").value = row.cells[9].innerText;
            document.getElementById("editGalaxy").value = row.cells[10].innerText;
            document.getElementById("editAccuracy").value = row.cells[11].innerText;

            modal.style.display = "block";
        });
    });

    // Закрытие модального окна
    span.onclick = function() {
        modal.style.display = "none";
    }

    window.onclick = function(event) {
        if (event.target == modal) {
            modal.style.display = "none";
        }
    }

    // Обработка кнопок удаления
    const deleteButtons = document.getElementsByClassName("delete-btn");
    Array.from(deleteButtons).forEach(button => {
        button.addEventListener("click", function() {
            const id = this.getAttribute("data-id");
            if (confirm("Вы уверены, что хотите удалить эту запись?")) {
                fetch(`/delete?id=${id}`, {
                    method: 'GET'
                }).then(response => {
                    if (response.ok) {
                        window.location.reload();
                    } else {
                        alert("Ошибка при удалении записи.");
                    }
                });
            }
        });
    });
});
