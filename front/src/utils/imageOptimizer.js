/**
 * 图片优化工具函数
 * 支持阿里云OSS图片处理参数
 */

/**
 * 优化图片URL，添加OSS图片处理参数
 * @param {string} url - 原始图片URL
 * @param {Object} options - 优化选项
 * @param {number} options.width - 目标宽度
 * @param {number} options.quality - 图片质量 (1-100)
 * @param {string} options.format - 目标格式 (webp, jpg, png)
 * @returns {string} 优化后的图片URL
 */
export function optimizeImage(url, options = {}) {
  if (!url || typeof url !== 'string') return url;

  // 默认配置
  const defaultOptions = {
    width: null,
    quality: 80,
    format: 'webp'
  };

  const config = { ...defaultOptions, ...options };

  // 检查是否是支持OSS处理的域名
  const supportedDomains = ['img.heliar.top', 'oss-cn-hongkong.aliyuncs.com', 'gvbresource.oss-cn-hongkong.aliyuncs.com'];
  const isOSSImage = supportedDomains.some(domain => url.includes(domain));

  if (!isOSSImage) return url;

  // 构建OSS图片处理参数
  const params = [];

  // 添加缩放参数
  if (config.width) {
    params.push(`image/resize,w_${config.width}`);
  }

  // 添加质量参数
  if (config.quality && config.quality !== 100) {
    params.push(`quality,q_${config.quality}`);
  }

  // 添加格式转换参数
  if (config.format) {
    params.push(`format,${config.format}`);
  }

  // 如果没有任何参数，返回原URL
  if (params.length === 0) return url;

  // 拼接参数
  const processParam = `?x-oss-process=${params.join('/')}`;

  // 如果URL已经有参数，需要特殊处理
  if (url.includes('?')) {
    return url.split('?')[0] + processParam;
  }

  return url + processParam;
}

/**
 * 预设的图片优化配置
 */
export const ImagePresets = {
  // 缩略图 (小图标、头像等)
  thumbnail: {
    width: 200,
    quality: 75,
    format: 'webp'
  },

  // 卡片图片 (文章卡片、成就卡片等)
  card: {
    width: 700,
    quality: 80,
    format: 'webp'
  },

  // 横幅图片 (Banner、背景图等)
  banner: {
    width: 1920,
    quality: 80,
    format: 'webp'
  },

  // 全尺寸 (文章详情图等)
  full: {
    width: 1200,
    quality: 85,
    format: 'webp'
  }
};

/**
 * 使用预设配置优化图片
 * @param {string} url - 原始图片URL
 * @param {string} preset - 预设名称 (thumbnail, card, banner, full)
 * @returns {string} 优化后的图片URL
 */
export function optimizeImageWithPreset(url, preset = 'card') {
  const config = ImagePresets[preset];
  if (!config) {
    console.warn(`Unknown preset: ${preset}, using default 'card' preset`);
    return optimizeImage(url, ImagePresets.card);
  }
  return optimizeImage(url, config);
}

/**
 * 生成响应式图片srcset
 * @param {string} url - 原始图片URL
 * @param {Array<number>} widths - 宽度数组
 * @returns {string} srcset字符串
 */
export function generateSrcset(url, widths = [400, 800, 1200, 1920]) {
  if (!url) return '';

  return widths
    .map(width => {
      const optimizedUrl = optimizeImage(url, { width, quality: 80, format: 'webp' });
      return `${optimizedUrl} ${width}w`;
    })
    .join(', ');
}

/**
 * 预加载图片
 * @param {string} url - 图片URL
 * @returns {Promise<void>}
 */
export function preloadImage(url) {
  return new Promise((resolve, reject) => {
    const img = new Image();
    img.onload = () => resolve();
    img.onerror = reject;
    img.src = url;
  });
}

/**
 * 批量预加载图片
 * @param {Array<string>} urls - 图片URL数组
 * @returns {Promise<Array>}
 */
export function preloadImages(urls) {
  return Promise.all(urls.map(url => preloadImage(url)));
}
