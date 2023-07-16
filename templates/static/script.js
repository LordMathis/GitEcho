// Retrieve backup repositories and populate the list
fetch('/api/v1/repository')
  .then(response => response.json())
  .then(data => {
    const backupRepoList = document.getElementById('backup-repo-list');
    data.forEach(backupRepo => {
      const li = document.createElement('li');
      li.className = 'backup-repo-item';
      li.textContent = backupRepo.name;

      li.addEventListener('click', () => {
        // Call a function to fill out the form with the selected backup repo's data
        populateFormWithBackupRepoData(backupRepo);
      });

      backupRepoList.appendChild(li);
    });
  });

  function populateFormWithBackupRepoData(backupRepo) {
    // Fill out the form fields with the backup repo's data
    document.getElementById('name').value = backupRepo.name;
    document.getElementById('remote-url').value = backupRepo.remote_url;
    document.getElementById('pull-interval').value = backupRepo.pull_interval;
    document.getElementById('git-username').value = backupRepo.credentials.git_username;
    document.getElementById('git-password').value = backupRepo.credentials.git_password;
    document.getElementById('git-key-path').value = backupRepo.credentials.git_key_path;

    // Clear the existing storage options
    const storageOptionsDiv = document.getElementById('storage-options');
    storageOptionsDiv.innerHTML = '';

    // Iterate through the backup repo's storages and add them to the form
    for (const storageName in backupRepo.storage) {
      const storage = backupRepo.storage[storageName];
      const storageForm = createStorageForm(storage);
      storageOptionsDiv.appendChild(storageForm);
    }
  }

// Show storage options form based on the selected storage type
function showStorageOptions() {
  const storageType = document.getElementById('storage-type').value;
  const storageOptionsDiv = document.getElementById('storage-options');

  if (storageType === 's3') {
    let storage = {
      type: 's3',
    }
    const storageForm = createStorageForm(storage);
    storageOptionsDiv.appendChild(storageForm);
  }
}

// Show storage options based on the selected storage type or add a remote storage form
// Add event listener for "Add Remote Storage" button
const addStorageBtn = document.getElementById('add-storage-btn');
addStorageBtn.addEventListener('click', () => {
  showStorageOptions(); // Show storage options form
});

// Add event listener for "Delete Backup Repo" button
const deleteBackupBtn = document.getElementById('delete-backup-btn');
deleteBackupBtn.addEventListener('click', () => {
  const name = document.getElementById('name').value;
  deleteBackupRepo(name);
});

// Add event listener for "Create Backup Repo" button
const createBackupBtn = document.getElementById('create-backup-btn');
createBackupBtn.addEventListener('click', (event) => {
  event.preventDefault(); // Prevent form submission

  // Get all storage forms
  const storageForms = document.querySelectorAll('.storage-form');


  // Retrieve form inputs
  const name = document.getElementById('name').value;
  const remoteUrl = document.getElementById('remote-url').value;
  const pullInterval = parseInt(document.getElementById('pull-interval').value);
  const gitUsername = document.getElementById('git-username').value;
  const gitPassword = document.getElementById('git-password').value;
  const gitKeyPath = document.getElementById('git-key-path').value;


  // Prepare data object for the API request
  const data = {
    name: name,
    remote_url: remoteUrl,
    pull_interval: pullInterval,
    credentials: {
      git_username: gitUsername,
      git_password: gitPassword,
      git_key_path: gitKeyPath,
    },
    storage: {},
  };

  // Iterate through each storage form and add the data to the storages object
  storageForms.forEach((form, index) => {
    const storageName = form.querySelector('input[name="storage-name"]').value;
    const storageEndpoint = form.querySelector('input[name="storage-endpoint"]').value;
    const storageRegion = form.querySelector('input[name="storage-region"]').value;
    const storageAccessKey = form.querySelector('input[name="storage-access-key"]').value;
    const storageSecretKey = form.querySelector('input[name="storage-secret-key"]').value;
    const storageBucketName = form.querySelector('input[name="storage-bucket-name"]').value;

    data.storage[storageName] = {
      name: storageName,
      type: form.getAttribute('data-storage-type'),
      endpoint: storageEndpoint,
      region: storageRegion,
      access_key: storageAccessKey,
      secret_key: storageSecretKey,
      bucket_name: storageBucketName,
    };
  });

  // Send the data to the API endpoint using fetch or your preferred AJAX method
  // Replace "/api/v1/backupRepos" with your actual API endpoint
  fetch('/api/v1/repository', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  })
    .then(response => response.json())
    .then(result => {
      console.log(result);
      // Handle the API response as needed
    })
    .catch(error => {
      console.error(error);
      // Handle errors if any
    });
});

function createStorageForm(storage) {
  const storageForm = document.createElement('div');
  storageForm.className = 'storage-form';
  storageForm.setAttribute('data-storage-type', storage.type);

  storageForm.innerHTML = `
    <h3>Storage Options</h3>

    <div class="form-grid">
      <label for="storage-name" class="required-label">Name:</label>
      <input type="text" name="storage-name" value="${storage.name || ''}" required>

      <label for="storage-endpoint">Endpoint:</label>
      <input type="text" name="storage-endpoint" value="${storage.endpoint || ''}">

      <label for="storage-region">Region:</label>
      <input type="text" name="storage-region" value="${storage.region || ''}">

      <label for="storage-access-key">Access Key:</label>
      <input type="text" name="storage-access-key" value="${storage.access_key || ''}">

      <label for="storage-secret-key">Secret Key:</label>
      <input type="text" name="storage-secret-key" value="${storage.secret_key || ''}">

      <label for="storage-bucket-name" class="required-label">Bucket Name:</label>
      <input type="text" name="storage-bucket-name" value="${storage.bucket_name || ''}" required>
    </div>
  `;

  return storageForm;
}

// Function to delete a backup repo
function deleteBackupRepo(repoName) {
  // Send a DELETE request to the API endpoint for deleting the backup repo
  fetch(`/api/v1/repository/${repoName}`, {
    method: 'DELETE',
  })
    .then(response => response.json())
    .then(result => {
      console.log(result);
      // Handle the API response as needed
    })
    .catch(error => {
      console.error(error);
      // Handle errors if any
    });
}