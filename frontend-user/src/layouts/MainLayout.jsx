import React from 'react'
import { Layout, Menu, Dropdown, Avatar, Button, Badge } from 'antd'
import {
  HomeOutlined,
  BookOutlined,
  FileTextOutlined,
  MessageOutlined,
  NotificationOutlined,
  UserOutlined,
  SettingOutlined,
  HeartOutlined,
  LogoutOutlined,
} from '@ant-design/icons'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'

const { Header, Content, Footer, Sider } = Layout

function MainLayout() {
  const navigate = useNavigate()
  const location = useLocation()
  const { user, logout, isAuthenticated } = useAuth()

  const menuItems = [
    { key: '/', icon: <HomeOutlined />, label: '首页' },
    { key: '/courses', icon: <BookOutlined />, label: '课程中心' },
    { key: '/exams', icon: <FileTextOutlined />, label: '在线考试' },
    { key: '/forums', icon: <MessageOutlined />, label: '论坛交流' },
    { key: '/announcements', icon: <NotificationOutlined />, label: '公告信息' },
  ]

  const handleMenuClick = ({ key }) => {
    navigate(key)
  }

  const userMenuItems = isAuthenticated
    ? [
        { key: '/profile', icon: <UserOutlined />, label: '个人信息' },
        { key: '/favorites', icon: <HeartOutlined />, label: '我的收藏' },
        { key: '/exam-records', icon: <FileTextOutlined />, label: '考试记录' },
        { type: 'divider' },
        {
          key: 'logout',
          icon: <LogoutOutlined />,
          label: '退出登录',
          onClick: () => {
            logout()
            navigate('/login')
          },
        },
      ]
    : []

  const selectedKey = location.pathname

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ display: 'flex', alignItems: 'center', padding: '0 24px', background: '#001529' }}>
        <div
          style={{ color: '#fff', fontSize: '20px', fontWeight: 'bold', cursor: 'pointer' }}
          onClick={() => navigate('/')}
        >
          📚 学习平台
        </div>

        <Menu
          theme="dark"
          mode="horizontal"
          selectedKeys={[selectedKey]}
          items={menuItems}
          onClick={handleMenuClick}
          style={{ flex: 1, justifyContent: 'center', borderBottom: 'none' }}
        />

        <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
          {isAuthenticated ? (
            <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
              <div style={{ cursor: 'pointer', display: 'flex', alignItems: 'center', gap: '8px' }}>
                <Avatar size={32} icon={<UserOutlined />} src={user?.avatar} />
                <span style={{ color: '#fff' }}>{user?.nickname || user?.username}</span>
              </div>
            </Dropdown>
          ) : (
            <>
              <Button type="link" onClick={() => navigate('/login')} style={{ color: '#fff' }}>
                登录
              </Button>
              <Button type="primary" onClick={() => navigate('/register')}>
                注册
              </Button>
            </>
          )}
        </div>
      </Header>

      <Content style={{ padding: '24px', background: '#f0f2f5' }}>
        <Outlet />
      </Content>

      <Footer style={{ textAlign: 'center', background: '#001529', color: '#fff' }}>
        学习平台 ©{new Date().getFullYear()} Created with React + Golang
      </Footer>
    </Layout>
  )
}

export default MainLayout
