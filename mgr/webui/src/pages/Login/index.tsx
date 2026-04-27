import {
  EyeInvisibleOutlined,
  LockOutlined,
  UserOutlined,
} from "@ant-design/icons";
import { Button, Checkbox, Form, Input } from "antd";

import styles from "./LoginPage.module.css";

type FeatureTone = "green" | "blue";

interface Feature {
  title: string;
  desc: string[];
  tone: FeatureTone;
  icon: "shield" | "bolt" | "pie" | "layers";
}

const features: Feature[] = [
  {
    title: "安全可靠",
    desc: ["多重加密保障", "数据安全"],
    tone: "green",
    icon: "shield",
  },
  {
    title: "高效稳定",
    desc: ["系统稳定高效", "处理海量数据"],
    tone: "blue",
    icon: "bolt",
  },
  {
    title: "精准计费",
    desc: ["灵活计费策略", "精准账单生成"],
    tone: "green",
    icon: "pie",
  },
  {
    title: "智能分析",
    desc: ["多维数据分析", "助力业务决策"],
    tone: "blue",
    icon: "layers",
  },
];

function AppMark() {
  return (
    <span className={styles.appMark} aria-hidden="true">
      <svg viewBox="0 0 48 48">
        <defs>
          <linearGradient id="markGradient" x1="6" x2="42" y1="4" y2="44">
            <stop stopColor="#27d59c" />
            <stop offset="1" stopColor="#1cbd82" />
          </linearGradient>
        </defs>
        <rect width="48" height="48" rx="10" fill="url(#markGradient)" />
        <path
          d="M15.5 13.5h16.2a2.8 2.8 0 0 1 2.8 2.8v12.2l-3.8-3.1-4.2 4.4 6.1 5.7H15.5a2.8 2.8 0 0 1-2.8-2.8V16.3a2.8 2.8 0 0 1 2.8-2.8Z"
          fill="#fff"
        />
        <path
          d="M19 19h9.6M19 23.5h9.6M19 28h6.2"
          stroke="#1cbd82"
          strokeLinecap="round"
          strokeWidth="2.4"
        />
        <path
          d="m30.2 28.5 5.6 5.5 3.7-4"
          fill="none"
          stroke="#fff"
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth="2.6"
        />
      </svg>
    </span>
  );
}

function FeatureIcon({ feature }: { feature: Feature }) {
  return (
    <span className={`${styles.featureIcon} ${styles[feature.tone]}`}>
      {feature.icon === "shield" ? (
        <svg viewBox="0 0 56 56">
          <path d="M28 8.5 43.5 14v12.5c0 9.7-6.4 18.4-15.5 21-9.1-2.6-15.5-11.3-15.5-21V14L28 8.5Z" fill="currentColor" />
          <path d="m23.5 28 3.4 3.4 6.6-7.2" fill="none" stroke="#fff" strokeLinecap="round" strokeLinejoin="round" strokeWidth="4" />
        </svg>
      ) : null}
      {feature.icon === "bolt" ? (
        <svg viewBox="0 0 56 56">
          <path d="M31.8 6 18.2 29.8h10.2L24.2 50l14.6-26.2H28.1L31.8 6Z" fill="currentColor" />
        </svg>
      ) : null}
      {feature.icon === "pie" ? (
        <svg viewBox="0 0 56 56">
          <path d="M28 8a20 20 0 1 0 20 20H28V8Z" fill="currentColor" />
          <path d="M32 8.4v15.8h15.6A20.2 20.2 0 0 0 32 8.4Z" fill="#e5fbf1" />
        </svg>
      ) : null}
      {feature.icon === "layers" ? (
        <svg viewBox="0 0 56 56">
          <path d="m28 8.5 18 9.4-18 9.5-18-9.5 18-9.4Z" fill="currentColor" />
          <path d="m12.5 25 15.5 8.1L43.5 25M12.5 33.5 28 41.6l15.5-8.1" fill="none" stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="5" />
        </svg>
      ) : null}
    </span>
  );
}

function FeatureGrid() {
  return (
    <div className={styles.featureGrid}>
      {features.map((feature) => (
        <article className={styles.featureCard} key={feature.title}>
          <FeatureIcon feature={feature} />
          <h3>{feature.title}</h3>
          {feature.desc.map((line) => (
            <p key={line}>{line}</p>
          ))}
        </article>
      ))}
    </div>
  );
}

function BrandIntro() {
  return (
    <section className={styles.brandIntro}>
      <div className={styles.brandName}>
        <AppMark />
        <span>账单计费系统</span>
      </div>

      <div className={styles.headline}>
        <h1>
          智能计费 <span>高效管理</span>
        </h1>
        <p>为企业提供一站式账单管理与计费解决方案</p>
      </div>

      <FeatureGrid />
    </section>
  );
}

function HeroArtwork() {
  return (
    <section className={styles.artwork} aria-hidden="true">
      <span className={`${styles.orb} ${styles.orbOne}`} />
      <span className={`${styles.orb} ${styles.orbTwo}`} />
      <span className={`${styles.diamond} ${styles.diamondOne}`} />
      <span className={`${styles.diamond} ${styles.diamondTwo}`} />
      <div className={styles.dotCloud} />

      <svg className={styles.artworkSvg} viewBox="0 0 620 430">
        <defs>
          <linearGradient id="screenBlue" x1="136" x2="438" y1="96" y2="282">
            <stop stopColor="#89adff" />
            <stop offset="1" stopColor="#4d7cf2" />
          </linearGradient>
          <linearGradient id="paperFront" x1="390" x2="556" y1="138" y2="318">
            <stop stopColor="#fff" />
            <stop offset="1" stopColor="#edf6ff" />
          </linearGradient>
          <linearGradient id="paperFold" x1="526" x2="572" y1="130" y2="176">
            <stop stopColor="#4dd8e8" />
            <stop offset="1" stopColor="#87e7f1" />
          </linearGradient>
          <filter id="artShadow" x="-30%" y="-30%" width="170%" height="170%">
            <feDropShadow dx="0" dy="22" stdDeviation="16" floodColor="#5b89ea" floodOpacity=".2" />
          </filter>
        </defs>

        <ellipse cx="306" cy="330" rx="248" ry="44" fill="#edf5ff" stroke="#cfe4ff" opacity=".75" />
        <ellipse cx="310" cy="313" rx="190" ry="29" fill="none" stroke="#c8def7" opacity=".55" />

        <path d="M132 288h318l-64 38H70l62-38Z" fill="#dceaff" filter="url(#artShadow)" />
        <path d="m154 272 9-176c.3-5.2 4.7-9.2 10-8.9l258 14.2a9.4 9.4 0 0 1 8.8 9.9L432 278l-278-6Z" fill="url(#screenBlue)" filter="url(#artShadow)" />
        <path d="m175 112 242 12-5.8 130.2-242-9.2L175 112Z" fill="#f7faff" opacity=".96" />
        <path d="M160 284h273l-43 26H118l42-26Z" fill="#6c96ff" />
        <path d="M206 294h142l-18 9H188l18-9Z" fill="#416fe5" opacity=".7" />

        <circle cx="190" cy="129" r="7" fill="#7fa2f4" />
        <rect x="204" y="123" width="36" height="11" rx="5.5" fill="#bed0fb" />
        <rect x="190" y="154" width="104" height="13" rx="3" fill="#dce8ff" />
        <rect x="190" y="178" width="132" height="13" rx="3" fill="#e5eeff" />
        <path d="m198 220 37-34 35 42 39-37 47 52" fill="none" stroke="#76a2ff" strokeLinecap="round" strokeLinejoin="round" strokeWidth="12" />
        <circle cx="356" cy="177" r="34" fill="#bed1ff" opacity=".62" />
        <path d="M356 177h34a34 34 0 0 1-14 29Z" fill="#6b95ff" opacity=".78" />
        <rect x="344" y="226" width="19" height="40" fill="#7da3ff" />
        <rect x="371" y="204" width="19" height="62" fill="#4778ef" />

        <path d="M398 130h132l36 40-18 140a15 15 0 0 1-16.4 13l-142-12a15 15 0 0 1-13.6-16.5l16.6-151.4c.9-7.6 4.4-13.1 5.4-13.1Z" fill="url(#paperFront)" filter="url(#artShadow)" />
        <path d="m530 130 36 40-45-4 9-36Z" fill="url(#paperFold)" />
        <path d="M425 172h78M424 201h96M421 231h100M418 260h76" stroke="#dbe7ff" strokeLinecap="round" strokeWidth="12" />
        <text x="418" y="207" fill="#79a1ff" fontSize="42" fontWeight="800">¥</text>

        <circle cx="548" cy="256" r="38" fill="#2fc984" filter="url(#artShadow)" />
        <path d="m532 255 13 14 27-31" fill="none" stroke="#fff" strokeLinecap="round" strokeLinejoin="round" strokeWidth="9" />
      </svg>
    </section>
  );
}

function LoginPanel() {
  return (
    <section className={styles.loginPanel} aria-labelledby="login-title">
      <div className={styles.loginHeading}>
        <h2 id="login-title">欢迎登录</h2>
        <p>请输入您的账号和密码</p>
      </div>

      <Form className={styles.loginForm} layout="vertical" autoComplete="off">
        <Form.Item
          className={styles.formItem}
          name="account"
          rules={[{ required: true, message: "请输入账号/邮箱/手机号" }]}
        >
          <Input
            size="large"
            prefix={<UserOutlined />}
            placeholder="请输入账号/邮箱/手机号"
          />
        </Form.Item>

        <Form.Item
          className={styles.formItem}
          name="password"
          rules={[{ required: true, message: "请输入密码" }]}
        >
          <Input.Password
            size="large"
            prefix={<LockOutlined />}
            iconRender={() => <EyeInvisibleOutlined />}
            placeholder="请输入密码"
          />
        </Form.Item>

        <div className={styles.formMeta}>
          <Checkbox>记住我</Checkbox>
          <button type="button">忘记密码?</button>
        </div>

        <Button type="primary" htmlType="submit" block className={styles.loginButton}>
          登录
        </Button>
      </Form>
    </section>
  );
}

export function LoginPage() {
  return (
    <main className={styles.loginPage}>
      <div className={styles.loginLayout}>
        <BrandIntro />
        <HeroArtwork />
        <LoginPanel />
      </div>
    </main>
  );
}
