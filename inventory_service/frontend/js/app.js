// Укажите URL вашего inventory_service
const API_BASE_URL = 'http://localhost:8080';

let editMode = false;

document.addEventListener('DOMContentLoaded', () => {
    fetchProducts();
    document.getElementById('show-create-form')
        .addEventListener('click', showCreateForm);
    document.getElementById('cancel-button')
        .addEventListener('click', hideForm);
    document.getElementById('product-form')
        .addEventListener('submit', onFormSubmit);
    document.getElementById('search-id')
        .addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                searchProduct();
            }
        });
});

function fetchProducts() {
    fetch(`${API_BASE_URL}/products`)
        .then(res => res.json())
        .then(renderTable)
        .catch(console.error);
}

function renderTable(products) {
    const tbody = document.querySelector('#products-table tbody');
    tbody.innerHTML = '';
    products.forEach(p => {
        const tr = document.createElement('tr');
        tr.innerHTML = `
      <td>${p.id}</td>
      <td>${p.name}</td>
      <td>${p.description}</td>
      <td>${p.price.toFixed(2)}</td>
      <td>
        <button onclick="editProduct(${p.id})">Изменить</button>
        <button onclick="deleteProduct(${p.id})">Удалить</button>
      </td>`;
        tbody.appendChild(tr);
    });
}

function showCreateForm() {
    editMode = false;
    document.getElementById('form-title').innerText = 'Создать продукт';
    document.getElementById('submit-button').innerText = 'Создать';
    document.getElementById('product-form').reset();
    document.getElementById('form-container').classList.remove('hidden');
}

function hideForm() {
    document.getElementById('form-container').classList.add('hidden');
}

function onFormSubmit(e) {
    e.preventDefault();
    const id = document.getElementById('product-id').value;
    const name = document.getElementById('product-name').value;
    const description = document.getElementById('product-description').value;
    const price = parseFloat(document.getElementById('product-price').value);

    const payload = { name, description, price };
    const method = editMode ? 'PATCH' : 'POST';
    const url = editMode
        ? `${API_BASE_URL}/products/${id}`
        : `${API_BASE_URL}/products`;

    fetch(url, {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
    })
        .then(res => {
            if (res.ok) {
                hideForm();
                fetchProducts();
            } else {
                return res.json().then(err => alert(err.error || 'Ошибка'));
            }
        })
        .catch(console.error);
}

function editProduct(id) {
    fetch(`${API_BASE_URL}/products/${id}`)
        .then(res => res.json())
        .then(p => {
            editMode = true;
            document.getElementById('form-title').innerText = 'Редактировать продукт';
            document.getElementById('submit-button').innerText = 'Сохранить';
            document.getElementById('product-id').value = p.id;
            document.getElementById('product-name').value = p.name;
            document.getElementById('product-description').value = p.description;
            document.getElementById('product-price').value = p.price;
            document.getElementById('form-container').classList.remove('hidden');
        })
        .catch(console.error);
}

function deleteProduct(id) {
    if (!confirm('Удалить продукт?')) return;
    fetch(`${API_BASE_URL}/products/${id}`, { method: 'DELETE' })
        .then(res => {
            if (res.ok) fetchProducts();
            else alert('Ошибка при удалении');
        })
        .catch(console.error);
}

function searchProduct() {
    const id = document.getElementById('search-id').value;
    if (!id) {
        alert('Введите ID для поиска');
        return;
    }

    fetch(`${API_BASE_URL}/products/${id}`)
        .then(res => {
            if (!res.ok) {
                if (res.status === 404) {
                    alert('Продукт не найден');
                    return null;
                }
                throw new Error('Ошибка поиска');
            }
            return res.json();
        })
        .then(product => {
            if (product) {
                renderTable([product]); // Отрисовываем как массив из одного элемента
            } else {
                renderTable([]); // Очищаем таблицу если продукт не найден
            }
        })
        .catch(error => {
            console.error('Ошибка:', error);
            alert('Произошла ошибка при поиске');
        });
}
