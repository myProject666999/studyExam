import React from 'react'
import { Layout, Menu, Avatar, Dropdown, Button } from 'antd'
import {
  DashboardOutlined,
  UserOutlined,
  BookOutlined,
  MessageOutlined,
  NotificationOutlined,
  PictureOutlined,
  FileTextOutlined,
  SettingOutlined,
  LogoutOutlined,
  DatabaseOutlined,
  OrderedListOutlined,
} from '@ant-design/icons'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'

const { Header, Sider, Content } = Layout

function MainLayout() {
  const navigate = useNavigate()
  const location = useLocation()
  const { admin, logout, isAuthenticated } = useAuth()

  React.useEffect(() => {
    if (!isAuthenticated) {
      navigate('/login')
    }
  }, [isAuthenticated, navigate])

  const menuItems = [
    { key: '/dashboard', icon: <DashboardOutlined />, label: '数据统计' },
    { key: '/users', icon: <UserOutlined />, label: '用户管理' },
    {
      key: 'course-group',
      icon: <BookOutlined />,
      label: '课程管理',
      children: [
        { key: '/course-types', label: '类型管理' },
        { key: '/courses', label: '课程列表' },
      ],
    },
    { key: '/forums', icon: <MessageOutlined />, label: '论坛管理' },
    { key: '/announcements', icon: <NotificationOutlined />, label: '公告管理' },
    { key: '/banners', icon: <PictureOutlined />, label: '轮播图管理' },
    {
      key: 'exam-group',
      icon: <FileTextOutlined />,
      label: '考试管理',
      children: [
        { key: '/exams', label: '考试列表' },
        { key: '/papers', label: '试卷管理' },
        { key: '/question-banks', label: '题库管理' },
        { key: '/questions', label: '试题管理' },
        { key: '/exam-records', label: '考试记录' },
      ],
    },
  ]

  const handleMenuClick = ({ key }) => {
    navigate(key)
  }

  const userMenuItems = [
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

  const selectedKey = location.pathname

  const getOpenKeys = () => {
    const path = location.pathname
    if (path.startsWith('/course-types') || path.startsWith('/courses')) {
      return ['course-group']
    }
    if (
      path.startsWith('/exams') ||
      path.startsWith('/papers') ||
      path.startsWith('/question-banks') ||
      path.startsWith('/questions') ||
      path.startsWith('/exam-records')
    ) {
      return ['exam-group']
    }
    return []
  }

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider width={220} theme="dark" collapsible>
        <div
          style={{
            height: 64,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            color: '#fff',
            fontSize: '18px',
            fontWeight: 'bold',
            background: 'rgba(255, 255, 255, 0.1)',
          }}
        >
          📚 管理后台
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[selectedKey]}
          defaultOpenKeys={getOpenKeys()}
          items={menuItems}
          onClick={handleMenuClick}
          style={{ borderRight: 0 }}
        />
      </Sider>

      <Layout>
        <Header
          style={{
            background: '#fff',
            padding: '0 24px',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'flex-end',
            boxShadow: '0 1px 4px rgba(0, 21, 41, 0.08)',
          }}
        >
          <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
            <div
              style={{ cursor: 'pointer', display: 'flex', alignItems: 'center', gap: '12px' }}
            >
              <Avatar size={36} icon={<UserOutlined />} src={admin?.avatar} />
              <span style={{ color: '#333' }}>
                {admin?.nickname || admin?.username || '管理员'}
              </span>
            </div>
          </Dropdown>
        </Header>

        <Content
          style={{
            margin: '24px',
            padding: '24px',
            background: '#fff',
            borderRadius: 8,
            minHeight: 'calc(100vh - 112px)',
          }}
        >
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}

export default MainLayout
