import React, { useEffect, useState } from 'react'
import {
  Row,
  Col,
  Card,
  Carousel,
  Typography,
  Tag,
  Button,
  Empty,
  Spin,
} from 'antd'
import {
  BookOutlined,
  StarOutlined,
  EyeOutlined,
  FireOutlined,
  HeartOutlined,
} from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import request from '../utils/request'
import { useAuth } from '../contexts/AuthContext'

const { Title, Text } = Typography

function Home() {
  const navigate = useNavigate()
  const { isAuthenticated } = useAuth()
  const [banners, setBanners] = useState([])
  const [recommendations, setRecommendations] = useState([])
  const [hotCourses, setHotCourses] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchData()
  }, [])

  const fetchData = async () => {
    setLoading(true)
    try {
      const [bannersRes, hotRes] = await Promise.all([
        request.get('/banners'),
        request.get('/hot-courses'),
      ])

      setBanners(bannersRes.data || [])
      setHotCourses(hotRes.data || [])

      if (isAuthenticated) {
        try {
          const recRes = await request.get('/recommendations')
          setRecommendations(recRes.data || [])
        } catch (error) {
          console.log('获取推荐失败')
        }
      }
    } catch (error) {
      console.error('获取首页数据失败:', error)
    } finally {
      setLoading(false)
    }
  }

  const CourseCard = ({ course, score }) => {
    return (
      <Card
        hoverable
        style={{ marginBottom: 16 }}
        onClick={() => navigate(`/courses/${course.id}`)}
        cover={
          <div
            style={{
              height: 160,
              background: `linear-gradient(135deg, #667eea 0%, #764ba2 100%)`,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              color: '#fff',
              fontSize: 48,
            }}
          >
            📚
          </div>
        }
      >
        <Card.Meta
          title={
            <div style={{ whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>
              {course.title}
            </div>
          }
          description={
            <div>
              <div style={{ marginBottom: 8 }}>
                <Tag color="blue">{course.course_type?.name}</Tag>
                <Text type="secondary" style={{ marginLeft: 8 }}>
                  {course.author}
                </Text>
              </div>
              <div style={{ display: 'flex', gap: 16, color: '#999' }}>
                <span>
                  <EyeOutlined /> {course.view_count}
                </span>
                <span>
                  <StarOutlined /> {course.collect_count}
                </span>
                {score && (
                  <span>
                    <HeartOutlined /> {score.toFixed(1)}分
                  </span>
                )}
              </div>
            </div>
          }
        />
      </Card>
    )
  }

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: 100 }}>
        <Spin size="large" />
      </div>
    )
  }

  return (
    <div>
      {banners.length > 0 && (
        <Carousel autoplay style={{ marginBottom: 32, borderRadius: 8, overflow: 'hidden' }}>
          {banners.map((banner) => (
            <div key={banner.id}>
              <div
                style={{
                  height: 300,
                  background: `linear-gradient(135deg, #667eea 0%, #764ba2 100%)`,
                  display: 'flex',
                  flexDirection: 'column',
                  alignItems: 'center',
                  justifyContent: 'center',
                  color: '#fff',
                }}
              >
                <Title level={2} style={{ color: '#fff', marginBottom: 16 }}>
                  {banner.title || '欢迎来到学习平台'}
                </Title>
                <Text style={{ color: '#fff', fontSize: 16 }}>
                  探索优质课程，开启学习之旅
                </Text>
              </div>
            </div>
          ))}
        </Carousel>
      )}

      {isAuthenticated && recommendations.length > 0 && (
        <div style={{ marginBottom: 32 }}>
          <div style={{ display: 'flex', alignItems: 'center', marginBottom: 16 }}>
            <HeartOutlined style={{ fontSize: 24, color: '#1890ff', marginRight: 8 }} />
            <Title level={3} style={{ margin: 0 }}>
              为你推荐
            </Title>
            <Text type="secondary" style={{ marginLeft: 16 }}>
              基于协同过滤推荐算法
            </Text>
          </div>
          <Row gutter={[16, 16]}>
            {recommendations.slice(0, 4).map((item) => (
              <Col xs={24} sm={12} md={6} key={item.course.id}>
                <CourseCard course={item.course} score={item.score} />
              </Col>
            ))}
          </Row>
        </div>
      )}

      <div style={{ marginBottom: 32 }}>
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 16 }}>
          <div style={{ display: 'flex', alignItems: 'center' }}>
            <FireOutlined style={{ fontSize: 24, color: '#ff4d4f', marginRight: 8 }} />
            <Title level={3} style={{ margin: 0 }}>
              热门课程
            </Title>
          </div>
          <Button type="link" onClick={() => navigate('/courses')}>
            查看更多 →
          </Button>
        </div>
        <Row gutter={[16, 16]}>
          {hotCourses.length > 0 ? (
            hotCourses.slice(0, 8).map((item) => (
              <Col xs={24} sm={12} md={6} key={item.course?.id || item.id}>
                <CourseCard course={item.course || item} score={item.score} />
              </Col>
            ))
          ) : (
            <Col span={24}>
              <Empty description="暂无课程" />
            </Col>
          )}
        </Row>
      </div>

      <Row gutter={[24, 24]}>
        <Col xs={24} sm={8}>
          <Card
            hoverable
            style={{ textAlign: 'center' }}
            onClick={() => navigate('/courses')}
          >
            <BookOutlined style={{ fontSize: 48, color: '#1890ff' }} />
            <Title level={4} style={{ marginTop: 16 }}>
              课程中心
            </Title>
            <Text type="secondary">丰富的在线课程资源</Text>
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card
            hoverable
            style={{ textAlign: 'center' }}
            onClick={() => navigate('/exams')}
          >
            <FireOutlined style={{ fontSize: 48, color: '#ff4d4f' }} />
            <Title level={4} style={{ marginTop: 16 }}>
              在线考试
            </Title>
            <Text type="secondary">检验学习成果</Text>
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card
            hoverable
            style={{ textAlign: 'center' }}
            onClick={() => navigate('/forums')}
          >
            <HeartOutlined style={{ fontSize: 48, color: '#52c41a' }} />
            <Title level={4} style={{ marginTop: 16 }}>
              论坛交流
            </Title>
            <Text type="secondary">分享学习心得</Text>
          </Card>
        </Col>
      </Row>
    </div>
  )
}

export default Home
