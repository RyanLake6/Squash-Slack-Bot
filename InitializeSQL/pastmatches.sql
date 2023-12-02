/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

DROP TABLE IF EXISTS `pastmatches`;
CREATE TABLE `pastmatches` (
  `player1` varchar(255) DEFAULT NULL,
  `player2` varchar(255) DEFAULT NULL,
  `winner` int DEFAULT NULL,
  `player1PrevPos` int DEFAULT NULL,
  `player2PrevPos` int DEFAULT NULL,
  `date` date DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `pastmatches` (`player1`, `player2`, `winner`, `player1PrevPos`, `player2PrevPos`, `date`) VALUES
('ryan', 'charles', 2, 7, 6, '2023-01-21');
INSERT INTO `pastmatches` (`player1`, `player2`, `winner`, `player1PrevPos`, `player2PrevPos`, `date`) VALUES
('ian', 'alex', 1, 1, 2, '2023-01-24');
INSERT INTO `pastmatches` (`player1`, `player2`, `winner`, `player1PrevPos`, `player2PrevPos`, `date`) VALUES
('peter', 'vanel', 1, 12, 13, '2023-01-20');
INSERT INTO `pastmatches` (`player1`, `player2`, `winner`, `player1PrevPos`, `player2PrevPos`, `date`) VALUES
('harrison', 'peter', 1, 11, 12, '2023-01-25'),
('alex', 'eli', 1, 2, 3, '2023-01-26'),
('ian', 'alex', 1, 1, 2, '2023-01-26');

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;