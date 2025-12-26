/**
 * Certificate Download Component Example
 * 
 * This is an example Vue.js component for downloading certificates.
 * Adapt it to your frontend framework (React, Angular, etc.)
 */

<template>
  <div class="certificates-container">
    <h2>My Certificates</h2>

    <div v-if="loading" class="loading">
      <p>Loading certificates...</p>
    </div>

    <div v-else-if="certificates.length === 0" class="no-certificates">
      <p>No certificates earned yet. Complete a course to earn one!</p>
    </div>

    <div v-else class="certificates-list">
      <div v-for="cert in certificates" :key="cert.id" class="certificate-card">
        <div class="certificate-header">
          <h3>{{ cert.course_name }}</h3>
          <span class="certificate-number">#{{ cert.certificate_number }}</span>
        </div>

        <div class="certificate-details">
          <p>
            <strong>Issued:</strong>
            {{ formatDate(cert.created_at) }}
          </p>
          <p>
            <strong>Course ID:</strong>
            {{ cert.course_id }}
          </p>
        </div>

        <div class="certificate-actions">
          <button
            @click="downloadCertificate(cert.id)"
            :disabled="downloadingId === cert.id"
            class="btn-download"
          >
            <span v-if="downloadingId === cert.id">Downloading...</span>
            <span v-else>📥 Download PDF</span>
          </button>

          <button
            @click="viewCertificate(cert.id)"
            class="btn-view"
          >
            👁️ View
          </button>
        </div>

        <div v-if="downloadError === cert.id" class="error-message">
          Failed to download certificate. Please try again.
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const studentId = ref(null)
const certificates = ref([])
const loading = ref(false)
const downloadingId = ref(null)
const downloadError = ref(null)

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8000/api/v1'

// Get student ID from Keycloak or store
onMounted(() => {
  loadCertificates()
})

/**
 * Load certificates for current student
 */
const loadCertificates = async () => {
  loading.value = true
  try {
    // Get student ID from your auth system
    const response = await fetch(
      `${API_BASE_URL}/certificates/${studentId.value}`,
      {
        headers: {
          Authorization: `Bearer ${getAuthToken()}`,
        },
      }
    )

    if (!response.ok) {
      throw new Error('Failed to load certificates')
    }

    certificates.value = await response.json()
  } catch (error) {
    console.error('Error loading certificates:', error)
  } finally {
    loading.value = false
  }
}

/**
 * Download certificate PDF using presigned URL
 */
const downloadCertificate = async (certificateId) => {
  downloadingId.value = certificateId
  downloadError.value = null

  try {
    // Get presigned download URL
    const response = await fetch(
      `${API_BASE_URL}/certificates/${certificateId}/download?expiration=3600`,
      {
        headers: {
          Authorization: `Bearer ${getAuthToken()}`,
        },
      }
    )

    if (!response.ok) {
      throw new Error('Failed to get download URL')
    }

    const data = await response.json()
    const downloadUrl = data.download_url

    // Download using presigned URL (no auth needed)
    const link = document.createElement('a')
    link.href = downloadUrl
    link.download = `certificate-${certificateId}.pdf`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)

  } catch (error) {
    console.error('Error downloading certificate:', error)
    downloadError.value = certificateId
  } finally {
    downloadingId.value = null
  }
}

/**
 * View certificate in new window
 */
const viewCertificate = async (certificateId) => {
  try {
    const response = await fetch(
      `${API_BASE_URL}/certificates/${certificateId}/download?expiration=3600`,
      {
        headers: {
          Authorization: `Bearer ${getAuthToken()}`,
        },
      }
    )

    if (!response.ok) {
      throw new Error('Failed to get certificate URL')
    }

    const data = await response.json()
    window.open(data.download_url, '_blank')

  } catch (error) {
    console.error('Error viewing certificate:', error)
  }
}

/**
 * Format date for display
 */
const formatDate = (dateString) => {
  const options = { year: 'numeric', month: 'long', day: 'numeric' }
  return new Date(dateString).toLocaleDateString('en-US', options)
}

/**
 * Get auth token (adjust to your auth system)
 */
const getAuthToken = () => {
  // From localStorage
  return localStorage.getItem('auth_token') ||
         // From sessionStorage
         sessionStorage.getItem('auth_token') ||
         // From cookies (example)
         getCookie('access_token')
}

/**
 * Get cookie value by name
 */
const getCookie = (name) => {
  const value = `; ${document.cookie}`
  const parts = value.split(`; ${name}=`)
  if (parts.length === 2) return parts.pop().split(';').shift()
  return null
}
</script>

<style scoped>
.certificates-container {
  max-width: 900px;
  margin: 0 auto;
  padding: 2rem;
}

h2 {
  margin-bottom: 2rem;
  color: #333;
}

.loading {
  text-align: center;
  padding: 2rem;
  color: #666;
}

.no-certificates {
  background: #f0f0f0;
  padding: 2rem;
  border-radius: 8px;
  text-align: center;
  color: #666;
}

.certificates-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 2rem;
}

.certificate-card {
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 1.5rem;
  background: white;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s, box-shadow 0.2s;
}

.certificate-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.certificate-header {
  display: flex;
  justify-content: space-between;
  align-items: start;
  margin-bottom: 1rem;
}

.certificate-header h3 {
  margin: 0;
  color: #1F4788;
  font-size: 1.2rem;
}

.certificate-number {
  background: #f0f0f0;
  padding: 0.25rem 0.75rem;
  border-radius: 4px;
  font-size: 0.9rem;
  color: #666;
}

.certificate-details {
  margin-bottom: 1.5rem;
  color: #666;
  font-size: 0.95rem;
}

.certificate-details p {
  margin: 0.5rem 0;
}

.certificate-actions {
  display: flex;
  gap: 1rem;
}

.btn-download,
.btn-view {
  flex: 1;
  padding: 0.75rem 1rem;
  border: none;
  border-radius: 4px;
  font-size: 0.95rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn-download {
  background: #1F4788;
  color: white;
}

.btn-download:hover:not(:disabled) {
  background: #153560;
}

.btn-download:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-view {
  background: #f0f0f0;
  color: #333;
}

.btn-view:hover {
  background: #e0e0e0;
}

.error-message {
  margin-top: 1rem;
  padding: 0.75rem;
  background: #fee;
  border: 1px solid #fcc;
  border-radius: 4px;
  color: #c33;
  font-size: 0.9rem;
}
</style>

<!-- ========================== REACT VERSION ========================== -->

/**
 * React Component Version
 */

import React, { useState, useEffect } from 'react'
import './Certificates.css'

export function CertificatesComponent() {
  const [studentId, setStudentId] = useState(null)
  const [certificates, setCertificates] = useState([])
  const [loading, setLoading] = useState(false)
  const [downloadingId, setDownloadingId] = useState(null)
  const [downloadError, setDownloadError] = useState(null)

  const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8000/api/v1'

  useEffect(() => {
    loadCertificates()
  }, [studentId])

  const loadCertificates = async () => {
    if (!studentId) return

    setLoading(true)
    try {
      const response = await fetch(
        `${API_BASE_URL}/certificates/${studentId}`,
        {
          headers: {
            Authorization: `Bearer ${getAuthToken()}`,
          },
        }
      )

      if (!response.ok) {
        throw new Error('Failed to load certificates')
      }

      setCertificates(await response.json())
    } catch (error) {
      console.error('Error loading certificates:', error)
    } finally {
      setLoading(false)
    }
  }

  const downloadCertificate = async (certificateId) => {
    setDownloadingId(certificateId)
    setDownloadError(null)

    try {
      const response = await fetch(
        `${API_BASE_URL}/certificates/${certificateId}/download?expiration=3600`,
        {
          headers: {
            Authorization: `Bearer ${getAuthToken()}`,
          },
        }
      )

      if (!response.ok) {
        throw new Error('Failed to get download URL')
      }

      const data = await response.json()
      const link = document.createElement('a')
      link.href = data.download_url
      link.download = `certificate-${certificateId}.pdf`
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
    } catch (error) {
      console.error('Error downloading certificate:', error)
      setDownloadError(certificateId)
    } finally {
      setDownloadingId(null)
    }
  }

  const getAuthToken = () => {
    return localStorage.getItem('auth_token') ||
           sessionStorage.getItem('auth_token')
  }

  if (loading) return <div className="loading">Loading certificates...</div>

  if (certificates.length === 0) {
    return <div className="no-certificates">No certificates earned yet</div>
  }

  return (
    <div className="certificates-container">
      <h2>My Certificates</h2>
      <div className="certificates-list">
        {certificates.map((cert) => (
          <div key={cert.id} className="certificate-card">
            <div className="certificate-header">
              <h3>{cert.course_name}</h3>
              <span className="certificate-number">#{cert.certificate_number}</span>
            </div>
            <div className="certificate-details">
              <p><strong>Issued:</strong> {new Date(cert.created_at).toLocaleDateString()}</p>
            </div>
            <div className="certificate-actions">
              <button
                onClick={() => downloadCertificate(cert.id)}
                disabled={downloadingId === cert.id}
                className="btn-download"
              >
                {downloadingId === cert.id ? 'Downloading...' : '📥 Download PDF'}
              </button>
            </div>
            {downloadError === cert.id && (
              <div className="error-message">Failed to download certificate</div>
            )}
          </div>
        ))}
      </div>
    </div>
  )
}

export default CertificatesComponent
