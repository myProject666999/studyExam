import React from 'react'
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import MainLayout from './layouts/MainLayout'
import Login from './pages/Login'
import Dashboard from './pages/Dashboard'
import UserManagement from './pages/UserManagement'
import CourseManagement from './pages/CourseManagement'
import CourseTypeManagement from './pages/CourseTypeManagement'
import ForumManagement from './pages/ForumManagement'
import AnnouncementManagement from './pages/AnnouncementManagement'
import BannerManagement from './pages/BannerManagement'
import ExamManagement from './pages/ExamManagement'
import ExamRecordManagement from './pages/ExamRecordManagement'
import PaperManagement from './pages/PaperManagement'
import QuestionManagement from './pages/QuestionManagement'
import QuestionBankManagement from './pages/QuestionBankManagement'
import { AuthProvider } from './contexts/AuthContext'

function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/" element={<MainLayout />}>
            <Route index element={<Navigate to="/dashboard" replace />} />
            <Route path="dashboard" element={<Dashboard />} />
            <Route path="users" element={<UserManagement />} />
            <Route path="courses" element={<CourseManagement />} />
            <Route path="course-types" element={<CourseTypeManagement />} />
            <Route path="forums" element={<ForumManagement />} />
            <Route path="announcements" element={<AnnouncementManagement />} />
            <Route path="banners" element={<BannerManagement />} />
            <Route path="exams" element={<ExamManagement />} />
            <Route path="exam-records" element={<ExamRecordManagement />} />
            <Route path="papers" element={<PaperManagement />} />
            <Route path="questions" element={<QuestionManagement />} />
            <Route path="question-banks" element={<QuestionBankManagement />} />
            <Route path="*" element={<Navigate to="/dashboard" replace />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  )
}

export default App
