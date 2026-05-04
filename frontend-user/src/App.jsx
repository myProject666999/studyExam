import React from 'react'
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import MainLayout from './layouts/MainLayout'
import Home from './pages/Home'
import Login from './pages/Login'
import Register from './pages/Register'
import CourseList from './pages/CourseList'
import CourseDetail from './pages/CourseDetail'
import SectionDetail from './pages/SectionDetail'
import ExamList from './pages/ExamList'
import ExamDetail from './pages/ExamDetail'
import ForumList from './pages/ForumList'
import ForumDetail from './pages/ForumDetail'
import AnnouncementList from './pages/AnnouncementList'
import AnnouncementDetail from './pages/AnnouncementDetail'
import Profile from './pages/Profile'
import Favorites from './pages/Favorites'
import ExamRecords from './pages/ExamRecords'
import { AuthProvider } from './contexts/AuthContext'

function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/" element={<MainLayout />}>
            <Route index element={<Home />} />
            <Route path="courses" element={<CourseList />} />
            <Route path="courses/:id" element={<CourseDetail />} />
            <Route path="sections/:id" element={<SectionDetail />} />
            <Route path="exams" element={<ExamList />} />
            <Route path="exams/:id" element={<ExamDetail />} />
            <Route path="forums" element={<ForumList />} />
            <Route path="forums/:id" element={<ForumDetail />} />
            <Route path="announcements" element={<AnnouncementList />} />
            <Route path="announcements/:id" element={<AnnouncementDetail />} />
            <Route path="profile" element={<Profile />} />
            <Route path="favorites" element={<Favorites />} />
            <Route path="exam-records" element={<ExamRecords />} />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  )
}

export default App
