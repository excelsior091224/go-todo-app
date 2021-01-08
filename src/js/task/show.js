'use strict';

document.addEventListener('DOMContentLoaded', () => {
    const deleteBtns = document.querySelectorAll('.tasks__item-delete');
    const csrfToken = document.getElementsByName('csrf')[0].content;

    const deleteTask = id => {
        let statusCode;

        fetch(`/api/tasks/${id}`, {
            method: 'DELETE',
            headers: { 'X-CSRF-Token': csrfToken}
        })
          .then(res => {
              statusCode = res.status;
              return res.json();
          })
          .then(data => {
              console.log(JSON.stringify(data));
              if (statusCode == 200) {
                window.location.href = '/';
              }
          })
          .catch(err => console.error(err));
    };

    for (let elm of deleteBtns) {
        elm.addEventListener('click', event => {
            event.preventDefault();
            deleteTask(elm.dataset.id);
        });
    }
});